package servicetoggle

import (
	"context"

	"github.com/avrebarra/goggle/internal/module/serviceaccesslog"
	domaintoggle "github.com/avrebarra/goggle/internal/module/servicetoggle/domaintoggle"
	storagetoggle "github.com/avrebarra/goggle/internal/module/servicetoggle/storetoggle"
	"github.com/avrebarra/goggle/internal/utils"
	"github.com/avrebarra/goggle/utils/ctxsaga"
	"github.com/avrebarra/goggle/utils/validator"
	"github.com/pkg/errors"
)

var _ Service = (*ServiceStd)(nil)

// mockable:true
type StorageToggle storagetoggle.Storage
type ServiceAccessLog serviceaccesslog.Service

type ServiceConfig struct {
	ToggleStore      StorageToggle    `validate:"required"`
	AccessLogService ServiceAccessLog `validate:"required"`
}

type ServiceStd struct {
	ServiceConfig
}

func NewService(cfg ServiceConfig) (out *ServiceStd, err error) {
	if err = validator.Validate(&cfg); err != nil {
		return nil, err
	}
	out = &ServiceStd{ServiceConfig: cfg}
	return
}

func (s *ServiceStd) DoListToggles(ctx context.Context, in ParamsDoListToggles) (out []domaintoggle.ToggleWithDetail, total int64, err error) {
	if err = validator.Validate(in); err != nil {
		err = errors.Wrap(err, "bad params")
		return
	}

	resp, tot, err := s.ToggleStore.FetchPaged(ctx, storagetoggle.ParamsFetchPaged(in))
	if err != nil {
		err = errors.Wrap(err, "paged fetch failed")
		return
	}

	out = resp
	total = tot
	return
}

func (s *ServiceStd) DoListStrayToggles(ctx context.Context, in ParamsDoListStrayToggles) (out []domaintoggle.ToggleWithDetail, total int64, err error) {
	if err = validator.Validate(in); err != nil {
		err = errors.Wrap(err, "bad params")
		return
	}

	resp, tot, err := s.ToggleStore.ListHeadlessAccessPaged(ctx, storagetoggle.ParamsListHeadlessAccessPaged(in))
	if err != nil {
		err = errors.Wrap(err, "paged fetch failed")
		return
	}

	out = resp
	total = tot
	return
}

func (s *ServiceStd) DoGetToggle(ctx context.Context, id string) (out domaintoggle.ToggleWithDetail, err error) {
	resp, _, err := s.ToggleStore.FetchPaged(ctx, storagetoggle.ParamsFetchPaged{
		FilterIDs: []string{id},
		SkipTotal: true,
	})
	if err != nil {
		err = errors.Wrap(err, "data fetching failed")
		return
	}
	if len(resp) == 0 {
		err = errors.Wrapf(ErrNotFound, "cannot find by id: %s", id)
		return
	}

	out = resp[0]
	return
}

func (s *ServiceStd) DoCreateToggle(ctx context.Context, in domaintoggle.Toggle) (out domaintoggle.Toggle, err error) {
	if err = validator.Validate(in); err != nil {
		err = errors.Wrap(err, "bad params")
		return
	}

	data1, err := s.ToggleStore.FetchToggleStatByID(ctx, in.ID)
	if errors.Is(err, storagetoggle.ErrStoreNotFound) {
		err = nil // discard err
	}
	if err != nil {
		err = errors.Wrap(err, "data fetching failed")
		return
	}
	if data1.ID != "" {
		err = errors.Wrapf(ErrAlreadyExists, "resource exists: %s", in.ID)
		return
	}

	data2, err := s.ToggleStore.UpsertToggle(ctx, in)
	if err != nil {
		err = errors.Wrap(err, "insert failed")
		return
	}

	out = domaintoggle.Toggle{}
	_ = utils.MorphFrom(&out, data2, nil)
	return
}

func (s *ServiceStd) DoUpdateToggle(ctx context.Context, id string, in domaintoggle.Toggle) (out domaintoggle.Toggle, err error) {
	if err = validator.Validate(in); err != nil {
		err = errors.Wrap(err, "bad params")
		return
	}

	resp1, _, err := s.ToggleStore.FetchPaged(ctx, storagetoggle.ParamsFetchPaged{FilterIDs: []string{id}, SkipTotal: true})
	if err != nil {
		err = errors.Wrap(err, "data check failed")
		return
	}
	if len(resp1) == 0 {
		err = errors.Wrapf(ErrNotFound, "resource not found: %s", in.ID)
		return
	}

	in.ID = id
	data, err := s.ToggleStore.UpsertToggle(ctx, in)
	if err != nil {
		err = errors.Wrap(err, "update failed")
		return
	}

	out = domaintoggle.Toggle{}
	_ = utils.MorphFrom(&out, data, nil)
	return
}

func (s *ServiceStd) DoRemoveToggle(ctx context.Context, id string) (out domaintoggle.Toggle, err error) {
	saga := ctxsaga.CreateSaga(ctx)

	resp1, _, err := s.ToggleStore.FetchPaged(ctx, storagetoggle.ParamsFetchPaged{FilterIDs: []string{id}, SkipTotal: true})
	if err != nil {
		err = errors.Wrap(err, "data check failed")
		return
	}
	if len(resp1) == 0 {
		err = errors.Wrapf(ErrNotFound, "resource not found: %s", id)
		return
	}
	data := resp1[0]

	err = func() (err error) {
		if err = s.ToggleStore.RemoveTogglesByIDs(ctx, []string{id}); err != nil {
			err = errors.Wrap(err, "removal failed")
			return
		}
		if err = s.AccessLogService.DeleteAccessLogByToggleID(ctx, id); err != nil {
			err = errors.Wrap(err, "access log removal failed")
			return
		}
		return
	}()
	if err != nil {
		if errRb := saga.Rollback(); errRb != nil {
			err = errors.Wrapf(err, "rollback failed: %s", errRb.Error())
		}
		return
	}

	if err = saga.Commit(); err != nil {
		err = errors.Wrap(err, "atomic ops committing failed")
		return
	}

	out = domaintoggle.Toggle{}
	_ = utils.MorphFrom(&out, data, nil)
	return
}

func (s *ServiceStd) DoStatToggle(ctx context.Context, id string) (out domaintoggle.ToggleStat, err error) {
	_ = s.AccessLogService.AddAccessLog(ctx, id)

	resp, err := s.ToggleStore.FetchToggleStatByID(ctx, id)
	if errors.Is(err, storagetoggle.ErrStoreNotFound) {
		err = errors.Wrapf(ErrNotFound, "resource not found: %s", id)
		return
	}
	if err != nil {
		err = errors.Wrap(err, "data fetching failed")
		return
	}

	out = resp
	return
}
