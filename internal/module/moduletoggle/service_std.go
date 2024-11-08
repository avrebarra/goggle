package moduletoggle

import (
	"context"
	"errors"
	"fmt"

	"github.com/avrebarra/goggle/internal/utils"
	"github.com/avrebarra/goggle/utils/validator"
)

var _ Service = (*ServiceStd)(nil)

type ServiceConfig struct {
	ToggleStore Store `validate:"required"`
}

type ServiceStd struct {
	Config ServiceConfig
}

func NewService(cfg ServiceConfig) (out *ServiceStd, err error) {
	if err = validator.Validate(&cfg); err != nil {
		return nil, err
	}
	out = &ServiceStd{Config: cfg}
	return
}

func (s *ServiceStd) DoListToggles(ctx context.Context, in ParamsDoListToggles) (out []ToggleWithDetail, total int64, err error) {
	if err = validator.Validate(in); err != nil {
		err = fmt.Errorf("bad params: %w", err)
		return
	}

	resp, tot, err := s.Config.ToggleStore.FetchPaged(ctx, ParamsFetchPaged(in))
	if err != nil {
		err = fmt.Errorf("paged fetch failed: %w", err)
		return
	}

	out = resp
	total = tot
	return
}

func (s *ServiceStd) DoListStrayToggles(ctx context.Context, in ParamsDoListStrayToggles) (out []ToggleWithDetail, total int64, err error) {
	if err = validator.Validate(in); err != nil {
		err = fmt.Errorf("bad params: %w", err)
		return
	}

	resp, tot, err := s.Config.ToggleStore.ListHeadlessAccessPaged(ctx, ParamsListHeadlessAccessPaged(in))
	if err != nil {
		err = fmt.Errorf("paged fetch failed: %w", err)
		return
	}

	out = resp
	total = tot
	return
}

func (s *ServiceStd) DoGetToggle(ctx context.Context, id string) (out ToggleWithDetail, err error) {
	resp, _, err := s.Config.ToggleStore.FetchPaged(ctx, ParamsFetchPaged{
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

func (s *ServiceStd) DoCreateToggle(ctx context.Context, in Toggle) (out Toggle, err error) {
	err = fmt.Errorf("not implemented")
	return
}

func (s *ServiceStd) DoUpdateToggle(ctx context.Context, id string) (out Toggle, err error) {
	err = fmt.Errorf("not implemented")
	return
}

func (s *ServiceStd) DoRemoveToggle(ctx context.Context, id string) (out Toggle, err error) {
	resp1, _, err := s.Config.ToggleStore.FetchPaged(ctx, ParamsFetchPaged{FilterIDs: []string{id}, SkipTotal: true})
	if err != nil {
		err = fmt.Errorf("data check failed: %w", err)
		return
	}
	if len(resp1) == 0 {
		err = fmt.Errorf("%w: %s", ErrNotFound, id)
		return
	}
	data := resp1[0]

	if err = s.Config.ToggleStore.RemoveTogglesByIDs(ctx, []string{id}); err != nil {
		err = fmt.Errorf("removal failed: %w", err)
		return
	}

	out = Toggle{}
	_ = utils.Translate(&out, data, nil)
	return
}

func (s *ServiceStd) DoStatToggle(ctx context.Context, id string) (out ToggleStat, err error) {
	resp, err := s.Config.ToggleStore.FetchToggleStatByID(ctx, id)
	if errors.Is(err, ErrStoreNotFound) {
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
