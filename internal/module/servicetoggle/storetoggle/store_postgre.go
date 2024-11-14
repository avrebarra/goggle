package storetoggle

import (
	"context"
	"time"

	"github.com/avrebarra/goggle/internal/module/servicetoggle/domaintoggle"
	"github.com/avrebarra/goggle/internal/utils"
	"github.com/avrebarra/goggle/utils/ctxsaga"
	"github.com/avrebarra/goggle/utils/validator"
	"github.com/guregu/null/v5"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	PostgreTableToggles    = "toggles"
	PostgreTableAccessLogs = "access_logs"
)

var _ Storage = (*StoragePostgre)(nil)

type ConfigStoragePostgre struct {
	DB *gorm.DB `validate:"required,structonly"`
}

type StoragePostgre struct {
	ConfigStoragePostgre
}

func NewStoragePostgre(cfg ConfigStoragePostgre) (out *StoragePostgre, err error) {
	if err = validator.Validate(&cfg); err != nil {
		err = errors.Wrap(err, "bad config")
		return
	}
	out = &StoragePostgre{ConfigStoragePostgre: cfg}
	return
}

func (s *StoragePostgre) FetchPaged(ctx context.Context, in ParamsFetchPaged) (out []domaintoggle.ToggleWithDetail, total int64, err error) {
	utils.ApplyDefaults(&in, &ParamsFetchPaged{Limit: 10, SortBy: "id", SortOrder: "asc"})
	if err = validator.Validate(in); err != nil {
		err = errors.Wrap(err, "bad params")
		return
	}

	type ResultData struct {
		ID               string      `gorm:"column:id"`
		Status           bool        `gorm:"column:status"`
		UpdatedAt        time.Time   `gorm:"column:updated_at"`
		LastAccessedAt   null.String `gorm:"column:last_accessed_at"`
		AccessFreqWeekly int         `gorm:"column:access_freq_weekly"`
	}

	// ***

	if err = validator.Validate(&in); err != nil {
		err = errors.Wrap(err, "bad params")
		return
	}

	data := []ResultData{}

	sqla := s.DB.Table(PostgreTableAccessLogs).
		Select("toggle_id, MAX(created_at) AS last_accessed_at").
		Group("toggle_id")
	sqaf := s.DB.Table(PostgreTableAccessLogs).
		Select("toggle_id, COUNT(*) AS access_count").
		Where("created_at > (?)", time.Now().AddDate(0, 0, -7)).
		Group("toggle_id")

	q := s.DB.Table(PostgreTableToggles+" as tog").
		Joins("LEFT JOIN (?) as la ON la.toggle_id = tog.id", sqla).
		Joins("LEFT JOIN (?) as af ON af.toggle_id = tog.id", sqaf).
		Select("id, status, updated_at, last_accessed_at, af.access_count as access_freq_weekly").
		Order(in.SortBy + " " + in.SortOrder)
	if in.FilterAccessed.Valid {
		if in.FilterAccessed.Bool {
			q = q.Where("last_accessed_at IS NOT NULL")
		} else {
			q = q.Where("last_accessed_at IS NULL")
		}
	}
	if in.FilterIDs != nil {
		q = q.Where("tog.id IN ?", in.FilterIDs)
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

	out = []domaintoggle.ToggleWithDetail{}
	for _, d := range data {
		val := domaintoggle.ToggleWithDetail{}
		if d.LastAccessedAt.Valid {
			t, _ := time.Parse(time.DateTime+"-07:00", d.LastAccessedAt.String)
			val.LastAccessedAt = t
		}
		utils.MorphFrom(&val, &d, nil)
		out = append(out, val)
	}

	return
}

func (s *StoragePostgre) ListHeadlessAccessPaged(ctx context.Context, in ParamsListHeadlessAccessPaged) (out []domaintoggle.ToggleWithDetail, total int64, err error) {
	utils.ApplyDefaults(&in, &ParamsListHeadlessAccessPaged{Limit: 10, SortBy: "id", SortOrder: "asc"})
	if err = validator.Validate(in); err != nil {
		err = errors.Wrap(err, "bad params")
		return
	}

	type ResultData struct {
		ID               string    `gorm:"column:toggle_id"`
		LastAccessedAt   time.Time `gorm:"column:last_accessed_at"`
		AccessFreqWeekly int       `gorm:"column:access_freq_weekly"`
	}

	// ***

	if err = validator.Validate(&in); err != nil {
		err = errors.Wrap(err, "bad params")
		return
	}

	data := []ResultData{}

	sqacs := s.DB.Table(PostgreTableAccessLogs).
		Group("toggle_id").
		Select("toggle_id, MAX(created_at) AS last_accessed_at")
	sqtog := s.DB.Table(PostgreTableToggles).
		Select("id, TRUE as hit")
	sqaf := s.DB.Table(PostgreTableAccessLogs).
		Select("toggle_id, COUNT(*) AS access_count").
		Where("created_at > ?", time.Now().AddDate(0, 0, -7)).
		Group("toggle_id")

	q := s.DB.Table("(?) as acs", sqacs).
		Joins("LEFT JOIN (?) as tog ON tog.id = acs.toggle_id", sqtog).
		Joins("LEFT JOIN (?) as af ON af.toggle_id = acs.toggle_id", sqaf).
		Where("tog.hit IS NULL").
		Select("acs.toggle_id, acs.last_accessed_at, af.access_count as access_freq_weekly").
		Order(in.SortBy + " " + in.SortOrder)
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

	for _, d := range data {
		val := domaintoggle.ToggleWithDetail{}
		utils.MorphFrom(&val, &d, nil)
		out = append(out, val)
	}

	return
}

func (s *StoragePostgre) FetchToggleStatByID(ctx context.Context, id string) (out domaintoggle.ToggleStat, err error) {
	type ResultData struct {
		ID     string `gorm:"column:id"`
		Status bool   `gorm:"column:status"`
	}

	// ***

	data := ResultData{}
	q := s.DB.Table(PostgreTableToggles+" as tog").
		Where("tog.id = ?", id).
		Select("tog.id, tog.status").
		First(&data)
	err = q.Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.Wrapf(ErrStoreNotFound, "id: %s", id)
		return
	}
	if err != nil {
		err = errors.Wrap(err, "db fetch failed")
		return
	}
	out = domaintoggle.ToggleStat(data)
	return
}

func (s *StoragePostgre) UpsertToggle(ctx context.Context, in domaintoggle.Toggle) (out domaintoggle.Toggle, err error) {
	type ParamData struct {
		ID        string    `gorm:"column:id" validate:"-"`
		Status    bool      `gorm:"column:status" validate:"-"`
		UpdatedAt time.Time `gorm:"column:updated_at" validate:"required"`
	}

	// ***

	data := ParamData(in)
	utils.ApplyDefaults(&data, &ParamData{UpdatedAt: time.Now(), Status: false})
	if err = validator.Validate(data); err != nil {
		err = errors.Wrap(err, "bad data")
		return
	}
	if data.ID != "" {
		data.UpdatedAt = time.Now()
	}

	execer := s.DB

	if saga, ok := ctxsaga.GetSaga(ctx); ok {
		execer = s.DB.Begin()
		defer func() {
			saga.AddRollbackFx(func() error { return execer.Rollback().Error })
			saga.AddCommitFx(func() error { return execer.Commit().Error })
		}()
	}

	q := execer.Table(PostgreTableToggles)
	if data.ID != "" {
		q = q.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{"status", "updated_at"}),
		})
	}
	q = q.Create(&in)
	err = q.Error
	if err != nil {
		return
	}

	out = in
	return
}

func (s *StoragePostgre) RemoveTogglesByIDs(ctx context.Context, ids []string) (err error) {
	execer := s.DB

	if saga, ok := ctxsaga.GetSaga(ctx); ok {
		execer = s.DB.Begin()
		defer func() {
			saga.AddRollbackFx(func() error { return execer.Rollback().Error })
			saga.AddCommitFx(func() error { return execer.Commit().Error })
		}()
	}

	type ResultData struct{}
	q := execer.Table(PostgreTableToggles+" as tog").
		Delete(&ResultData{}, "tog.id IN (?)", ids)
	err = q.Error
	if err != nil {
		return
	}

	return
}
