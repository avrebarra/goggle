package serviceaccesslog

import (
	"context"
	"fmt"
	"time"

	domainaccesslog "github.com/avrebarra/goggle/internal/module/serviceaccesslog/domain"
	storageaccesslog "github.com/avrebarra/goggle/internal/module/serviceaccesslog/storage"
	"github.com/avrebarra/goggle/utils/validator"
)

var _ Service = (*ServiceStd)(nil)

type ServiceConfig struct {
	AccessLogStore storageaccesslog.Storage `validate:"required"`
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

func (s *ServiceStd) DoListLogs(ctx context.Context, in ParamsDoListLogs) (out []domainaccesslog.AccessLog, total int64, err error) {
	if err = validator.Validate(in); err != nil {
		err = fmt.Errorf("bad params: %w", err)
		return
	}

	resp, tot, err := s.AccessLogStore.FetchPaged(ctx, storageaccesslog.ParamsFetchPaged(in))
	if err != nil {
		err = fmt.Errorf("paged fetch failed: %w", err)
		return
	}

	out = resp
	total = tot
	return
}

func (s *ServiceStd) AddAccessLog(ctx context.Context, toggleid string) (err error) {
	_, err = s.AccessLogStore.CreateLog(ctx, domainaccesslog.AccessLog{ToggleID: toggleid, CreatedAt: time.Now()})
	if err != nil {
		err = fmt.Errorf("insert failed: %w", err)
		return
	}
	return
}

func (s *ServiceStd) DeleteAccessLogByToggleID(ctx context.Context, toggleid string) (err error) {
	if err = s.AccessLogStore.DeleteAllByToggleIDs(ctx, []string{toggleid}); err != nil {
		err = fmt.Errorf("deletion failed: %w", err)
		return
	}
	return
}
