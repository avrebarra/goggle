package moduletoggle

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/avrebarra/goggle/internal/utils"
	"github.com/avrebarra/goggle/utils/validator"
	"github.com/guregu/null/v5"
	"gorm.io/gorm"
)

var (
	SQLiteTableToggles    = "toggles"
	SQLiteTableAccessLogs = "access_logs"
)

var _ Store = (*StoreSQLite)(nil)

type ConfigStoreSQLite struct {
	DB *gorm.DB `validate:"required,structonly"`
}

type StoreSQLite struct {
	ConfigStoreSQLite
}

func NewStoreSQLite(cfg ConfigStoreSQLite) (out *StoreSQLite, err error) {
	if err = validator.Validate(&cfg); err != nil {
		err = fmt.Errorf("bad config: %w", err)
		return
	}
	out = &StoreSQLite{ConfigStoreSQLite: cfg}
	return
}

func (s *StoreSQLite) FetchPaged(ctx context.Context, in ParamsFetchPaged) (out []ToggleWithDetail, total int64, err error) {
	utils.ApplyDefaults(&in, &ParamsFetchPaged{Limit: 10, SortBy: "id", SortOrder: "asc"})
	if err = validator.Validate(in); err != nil {
		err = fmt.Errorf("bad params: %w", err)
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
		err = fmt.Errorf("bad params: %w", err)
		return
	}

	data := []ResultData{}

	sqla := s.DB.Table(SQLiteTableAccessLogs).
		Select("toggle_id, MAX(created_at) AS last_accessed_at").
		Group("toggle_id")
	sqaf := s.DB.Table(SQLiteTableAccessLogs).
		Select("toggle_id, COUNT(*) AS access_count").
		Where("created_at > ?", time.Now().AddDate(0, 0, -7)).
		Group("toggle_id")

	q := s.DB.Table(SQLiteTableToggles+" as tog").
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
		err = fmt.Errorf("db fetch failed: %w", err)
		return
	}

	out = []ToggleWithDetail{}
	for _, d := range data {
		t, _ := time.Parse(time.DateTime, d.LastAccessedAt.String)

		val := ToggleWithDetail{}
		utils.Translate(&val, &d, &ToggleWithDetail{LastAccessedAt: t})

		out = append(out, val)
	}

	return
}

func (s *StoreSQLite) ListHeadlessAccessPaged(ctx context.Context, in ParamsListHeadlessAccessPaged) (out []ToggleWithDetail, total int64, err error) {
	utils.ApplyDefaults(&in, &ParamsListHeadlessAccessPaged{Limit: 10, SortBy: "id", SortOrder: "asc"})
	if err = validator.Validate(in); err != nil {
		err = fmt.Errorf("bad params: %w", err)
		return
	}

	type ResultData struct {
		ID               string `gorm:"column:toggle_id"`
		LastAccessedAt   string `gorm:"column:last_accessed_at"`
		AccessFreqWeekly int    `gorm:"column:access_freq_weekly"`
	}

	// ***

	if err = validator.Validate(&in); err != nil {
		err = fmt.Errorf("bad params: %w", err)
		return
	}

	data := []ResultData{}

	sqacs := s.DB.Table(SQLiteTableAccessLogs).
		Group("toggle_id").
		Select("toggle_id, MAX(created_at) AS last_accessed_at")
	sqtog := s.DB.Table(SQLiteTableToggles).
		Select("id, TRUE as hit")
	sqaf := s.DB.Table(SQLiteTableAccessLogs).
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
		err = fmt.Errorf("db fetch failed: %w", err)
		return
	}

	out = []ToggleWithDetail{}
	for _, d := range data {
		t, _ := time.Parse(time.DateTime, d.LastAccessedAt)

		val := ToggleWithDetail{}
		utils.Translate(&val, &d, &ToggleWithDetail{LastAccessedAt: t})

		out = append(out, val)
	}

	return
}

func (s *StoreSQLite) FetchToggleStatByID(ctx context.Context, id string) (out ToggleStat, err error) {
	type ResultData struct {
		ID     string `gorm:"column:id"`
		Status bool   `gorm:"column:status"`
	}

	// ***

	data := ResultData{}
	q := s.DB.Table(SQLiteTableToggles+" as tog").
		Where("tog.id = ?", id).
		Select("tog.id, tog.status").
		First(&data)
	err = q.Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = fmt.Errorf("%w: %s", ErrStoreNotFound, id)
		return
	}
	if err != nil {
		err = fmt.Errorf("db fetch failed: %w", err)
		return
	}
	jstr, _ := json.MarshalIndent(map[string]any{"data": data}, "", "  ")
	fmt.Println(string(jstr))

	out = ToggleStat(data)
	return
}

func (s *StoreSQLite) RemoveTogglesByIDs(ctx context.Context, ids []string) (err error) {
	type ResultData struct{}
	q := s.DB.Table(SQLiteTableToggles+" as tog").
		Delete(&ResultData{}, "tog.id IN (?)", ids)
	err = q.Error
	if err != nil {
		return
	}
	return
}
