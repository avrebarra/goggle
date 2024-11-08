package moduletoggle

import (
	"context"
	"fmt"

	"github.com/guregu/null/v5"
)

var (
	ErrNotFound = fmt.Errorf("not found")
)

type Service interface {
	DoListToggles(ctx context.Context, in ParamsDoListToggles) (out []ToggleWithDetail, total int64, err error)
	DoListStrayToggles(ctx context.Context, in ParamsDoListStrayToggles) (out []ToggleWithDetail, total int64, err error)
	DoGetToggle(ctx context.Context, id string) (out ToggleWithDetail, err error)
	DoCreateToggle(ctx context.Context, in Toggle) (out Toggle, err error)
	DoUpdateToggle(ctx context.Context, id string) (out Toggle, err error)
	DoRemoveToggle(ctx context.Context, id string) (out Toggle, err error)
	DoStatToggle(ctx context.Context, id string) (out ToggleStat, err error)
}

type ParamsDoListToggles struct {
	Offset    int
	Limit     int
	SortBy    string
	SortOrder string

	FilterIDs      []string
	FilterAccessed null.Bool

	SkipTotal bool
}

type ParamsDoListStrayToggles struct {
	Offset    int
	Limit     int
	SortBy    string
	SortOrder string

	SkipTotal bool
}
