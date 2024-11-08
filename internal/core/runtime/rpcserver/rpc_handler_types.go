package rpcserver

import (
	"time"

	"github.com/guregu/null/v5"
)

// request & responses

type ReqPing struct{}
type RespPing struct {
	Version   string    `json:"version"`
	StartedAt time.Time `json:"startedAt"`
	Uptime    string    `json:"uptime"`
}

type ReqListToggles struct {
	Offset         int       `json:"offset"`
	Limit          int       `json:"limit"`
	SortBy         string    `json:"sortBy"`
	SortOrder      string    `json:"sortOrder"`
	FilterAccessed null.Bool `json:"filterAccessed"`
}
type RespListToggles struct {
	Items []ToggleWithDetail `json:"items"`
	Total int64              `json:"total"`
}

type ReqListStrayToggles struct {
	Offset    int    `json:"offset"`
	Limit     int    `json:"limit"`
	SortBy    string `json:"sortBy"`
	SortOrder string `json:"sortOrder"`
}
type RespListStrayToggles struct {
	Items []ToggleStatLog `json:"items"`
	Total int64           `json:"total"`
}

type ReqGetToggle struct {
	ID string `json:"id"`
}

type ReqUpdateToggle struct {
	ID string `json:"id"`
}

type ReqRemoveToggle struct {
	ID string `json:"id"`
}

type ReqStatToggle struct {
	ID string `json:"id"`
}

// entities

type Toggle struct {
	ID        string    `json:"id" validate:"required"`
	Status    bool      `json:"status" validate:"-"`
	UpdatedAt time.Time `json:"updatedAt" validate:"required"`
}

type ToggleWithDetail struct {
	ID               string    `json:"id"`
	Status           bool      `json:"status"`
	UpdatedAt        time.Time `json:"updatedAt"`
	LastAccessedAt   null.Time `json:"lastAccessedAt"`
	AccessFreqWeekly int       `json:"accessFreqWeekly"`
}

type ToggleStatLog struct {
	ID               string    `json:"id"`
	LastAccessedAt   null.Time `json:"lastAccessedAt"`
	AccessFreqWeekly int       `json:"accessFreqWeekly"`
}

type ToggleCompact struct {
	ID     string `json:"id" validate:"required"`
	Status bool   `json:"status" validate:"-"`
}
