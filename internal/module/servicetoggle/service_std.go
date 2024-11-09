package servicetoggle

import (
	"context"
	"errors"
	"fmt"

	"github.com/avrebarra/goggle/internal/module/serviceaccesslog"
	domaintoggle "github.com/avrebarra/goggle/internal/module/servicetoggle/domain"
	storagetoggle "github.com/avrebarra/goggle/internal/module/servicetoggle/storage"

	"github.com/avrebarra/goggle/internal/utils"
	"github.com/avrebarra/goggle/utils/ctxsaga"
	"github.com/avrebarra/goggle/utils/validator"
)

var _ Service = (*ServiceStd)(nil)

type ServiceConfig struct {
	ToggleStore      storagetoggle.Storage    `validate:"required"`
	AccessLogService serviceaccesslog.Service `validate:"required"`
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
		err = fmt.Errorf("bad params: %w", err)
		return
	}

	resp, tot, err := s.ToggleStore.FetchPaged(ctx, storagetoggle.ParamsFetchPaged(in))
	if err != nil {
		err = fmt.Errorf("paged fetch failed: %w", err)
		return
	}

	out = resp
	total = tot
	return
}

func (s *ServiceStd) DoListStrayToggles(ctx context.Context, in ParamsDoListStrayToggles) (out []domaintoggle.ToggleWithDetail, total int64, err error) {
	if err = validator.Validate(in); err != nil {
		err = fmt.Errorf("bad params: %w", err)
		return
	}

	resp, tot, err := s.ToggleStore.ListHeadlessAccessPaged(ctx, storagetoggle.ParamsListHeadlessAccessPaged(in))
	if err != nil {
		err = fmt.Errorf("paged fetch failed: %w", err)
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
		err = fmt.Errorf("data fetching failed: %w", err)
		return
	}
	if len(resp) == 0 {
		err = fmt.Errorf("%w: %s", ErrNotFound, id)
		return
	}

	out = resp[0]
	return
}

func (s *ServiceStd) DoCreateToggle(ctx context.Context, in domaintoggle.Toggle) (out domaintoggle.Toggle, err error) {
	data1, err := s.ToggleStore.FetchToggleStatByID(ctx, in.ID)
	if errors.Is(err, storagetoggle.ErrStoreNotFound) {
		err = nil // discard err
	}
	if err != nil {
		err = fmt.Errorf("data fetching failed: %w", err)
		return
	}
	if data1.ID != "" {
		err = fmt.Errorf("%w: %s", ErrAlreadyExists, in.ID)
		return
	}

	data2, err := s.ToggleStore.UpsertToggle(ctx, in)
	if err != nil {
		err = fmt.Errorf("insert failed: %w", err)
		return
	}

	out = domaintoggle.Toggle{}
	_ = utils.MorphFrom(&out, data2, nil)
	return
}

func (s *ServiceStd) DoUpdateToggle(ctx context.Context, id string, in domaintoggle.Toggle) (out domaintoggle.Toggle, err error) {
	resp1, _, err := s.ToggleStore.FetchPaged(ctx, storagetoggle.ParamsFetchPaged{FilterIDs: []string{id}, SkipTotal: true})
	if err != nil {
		err = fmt.Errorf("data check failed: %w", err)
		return
	}
	if len(resp1) == 0 {
		err = fmt.Errorf("%w: %s", ErrNotFound, id)
		return
	}

	in.ID = id
	data, err := s.ToggleStore.UpsertToggle(ctx, in)
	if err != nil {
		err = fmt.Errorf("update failed: %w", err)
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
		err = fmt.Errorf("data check failed: %w", err)
		return
	}
	if len(resp1) == 0 {
		err = fmt.Errorf("%w: %s", ErrNotFound, id)
		return
	}
	data := resp1[0]

	err = func() (err error) {
		if err = s.ToggleStore.RemoveTogglesByIDs(ctx, []string{id}); err != nil {
			err = fmt.Errorf("removal failed: %w", err)
			return
		}
		if err = s.AccessLogService.DeleteAccessLogByToggleID(ctx, id); err != nil {
			err = fmt.Errorf("access log removal failed: %w", err)
			return
		}
		return
	}()
	if err != nil {
		if errRb := saga.Rollback(); errRb != nil {
			err = fmt.Errorf("%w: rollback failed: %v", err, errRb)
		}
		return
	}

	if err = saga.Commit(); err != nil {
		err = fmt.Errorf("atomic ops committing failed: %w", err)
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
		err = fmt.Errorf("%w: %s", ErrNotFound, id)
		return
	}
	if err != nil {
		err = fmt.Errorf("data fetching failed: %w", err)
		return
	}

	out = resp
	return
}
