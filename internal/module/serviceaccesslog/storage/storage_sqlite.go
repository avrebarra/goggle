package storage

import (
	"context"
	"time"

	domainaccesslog "github.com/avrebarra/goggle/internal/module/serviceaccesslog/domain"
	"github.com/avrebarra/goggle/internal/utils"
	"github.com/avrebarra/goggle/utils/ctxsaga"
	"github.com/avrebarra/goggle/utils/validator"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

var (
	SQLiteTableAccessLogs = "access_logs"
)

var _ Storage = (*StorageSQLite)(nil)

type ConfigStorageSQLite struct {
	DB *gorm.DB `validate:"required,structonly"`
}

type StorageSQLite struct {
	ConfigStorageSQLite
}

func NewStorageSQLite(cfg ConfigStorageSQLite) (out *StorageSQLite, err error) {
	if err = validator.Validate(&cfg); err != nil {
		err = errors.Wrap(err, "bad config")
		return
	}
	out = &StorageSQLite{ConfigStorageSQLite: cfg}
	return
}

func (s *StorageSQLite) FetchPaged(ctx context.Context, in ParamsFetchPaged) (out []domainaccesslog.AccessLog, total int64, err error) {
	utils.ApplyDefaults(&in, &ParamsFetchPaged{Limit: 10, SortBy: "id", SortOrder: "asc"})
	if err = validator.Validate(in); err != nil {
		err = errors.Wrap(err, "bad params")
		return
	}

	type ResultData struct {
		ID        int       `gorm:"column:id"`
		ToggleID  string    `gorm:"column:toggle_id"`
		CreatedAt time.Time `gorm:"column:created_at"`
	}

	// ***

	if err = validator.Validate(&in); err != nil {
		err = errors.Wrap(err, "bad params")
		return
	}

	data := []ResultData{}
	q := s.DB.Table(SQLiteTableAccessLogs + " as log").
		Select("id, toggle_id, created_at").
		Order(in.SortBy + " " + in.SortOrder)
	if in.FilterToggleIDs != nil {
		q = q.Where("log.toggle_id IN ?", in.FilterToggleIDs)
	}
	if !in.SkipTotal {
		q = q.Count(&total)
	}
	q = q.Limit(in.Limit).
		Offset(in.Offset).
		Find(&data)
	if err = q.Error; err != nil {
		err = errors.Wrap(err, "db fetch failed")
		return
	}

	out = []domainaccesslog.AccessLog{}
	for _, d := range data {
		val := domainaccesslog.AccessLog(d)
		out = append(out, val)
	}

	return
}

func (s *StorageSQLite) CreateLog(ctx context.Context, in domainaccesslog.AccessLog) (out domainaccesslog.AccessLog, err error) {
	type ParamData struct {
		ID        int       `gorm:"column:id" validate:"-"`
		ToggleID  string    `gorm:"column:toggle_id" validate:"required"`
		CreatedAt time.Time `gorm:"column:created_at" validate:"required"`
	}

	// ***

	data := ParamData(in)
	utils.ApplyDefaults(&data, &ParamData{CreatedAt: time.Now()})
	if err = validator.Validate(data); err != nil {
		err = errors.Wrap(err, "bad data")
		return
	}

	execer := s.DB

	if saga, ok := ctxsaga.GetSaga(ctx); ok {
		execer = s.DB.Begin()
		defer func() {
			saga.AddRollbackFx(func() error { return execer.Rollback().Error })
			saga.AddCommitFx(func() error { return execer.Commit().Error })
		}()
	}

	q := execer.Table(SQLiteTableAccessLogs).
		Create(&data)
	err = q.Error
	if err != nil {
		return
	}

	out = in
	return
}

func (s *StorageSQLite) DeleteAllByToggleIDs(ctx context.Context, toggleids []string) (err error) {
	execer := s.DB

	if saga, ok := ctxsaga.GetSaga(ctx); ok {
		execer = s.DB.Begin()
		defer func() {
			saga.AddRollbackFx(func() error { return execer.Rollback().Error })
			saga.AddCommitFx(func() error { return execer.Commit().Error })
		}()
	}

	type ResultData struct{}
	q := execer.Table(SQLiteTableAccessLogs+" as log").
		Delete(&ResultData{}, "log.toggle_id IN (?)", toggleids)
	err = q.Error
	if err != nil {
		return
	}

	return
}
