package servicetoggle

import (
	"context"

	domaintoggle "github.com/avrebarra/goggle/internal/module/servicetoggle/domaintoggle"
	"github.com/pkg/errors"

	"github.com/guregu/null/v5"
)

var (
	ErrNotFound      = errors.Errorf("not found")
	ErrAlreadyExists = errors.Errorf("already exists")
)

// mockable:true
type Service interface {
	DoListToggles(ctx context.Context, in ParamsDoListToggles) (out []domaintoggle.ToggleWithDetail, total int64, err error)
	DoListStrayToggles(ctx context.Context, in ParamsDoListStrayToggles) (out []domaintoggle.ToggleWithDetail, total int64, err error)
	DoGetToggle(ctx context.Context, id string) (out domaintoggle.ToggleWithDetail, err error)
	DoCreateToggle(ctx context.Context, in domaintoggle.Toggle) (out domaintoggle.Toggle, err error)
	DoUpdateToggle(ctx context.Context, id string, in domaintoggle.Toggle) (out domaintoggle.Toggle, err error)
	DoRemoveToggle(ctx context.Context, id string) (out domaintoggle.Toggle, err error)
	DoStatToggle(ctx context.Context, id string) (out domaintoggle.ToggleStat, err error)
}

type ParamsDoListToggles struct {
	Offset    int `validate:"min=0"`
	Limit     int `validate:"min=0,max=100"`
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
