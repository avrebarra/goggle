package storetoggle

import (
	"context"

	domaintoggle "github.com/avrebarra/goggle/internal/module/servicetoggle/domaintoggle"
	"github.com/guregu/null/v5"
	"github.com/pkg/errors"
)

var (
	ErrStoreNotFound = errors.New("not found")
)

type Storage interface {
	FetchPaged(ctx context.Context, in ParamsFetchPaged) (out []domaintoggle.ToggleWithDetail, total int64, err error)
	ListHeadlessAccessPaged(ctx context.Context, in ParamsListHeadlessAccessPaged) (out []domaintoggle.ToggleWithDetail, total int64, err error)
	FetchToggleStatByID(ctx context.Context, id string) (out domaintoggle.ToggleStat, err error)
	RemoveTogglesByIDs(ctx context.Context, ids []string) (err error)
	UpsertToggle(ctx context.Context, in domaintoggle.Toggle) (out domaintoggle.Toggle, err error)
}

type ParamsFetchPaged struct {
	Offset    int    `validate:"min=0"`
	Limit     int    `validate:"min=0,max=100"`
	SortBy    string `validate:"oneof=id name"`
	SortOrder string `validate:"oneof=asc desc"`

	FilterIDs      []string  `validate:"-"`
	FilterAccessed null.Bool `validate:"-"`

	SkipTotal bool `validate:"-"`
}

type ParamsListHeadlessAccessPaged struct {
	Offset    int    `validate:"min=0"`
	Limit     int    `validate:"min=0,max=100"`
	SortBy    string `validate:"oneof=id name"`
	SortOrder string `validate:"oneof=asc desc"`

	SkipTotal bool `validate:"-"`
}
