package moduletoggle

import (
	"context"
	"errors"

	"github.com/guregu/null/v5"
)

var (
	ErrStoreNotFound = errors.New("not found")
)

type Store interface {
	FetchPaged(ctx context.Context, in ParamsFetchPaged) (out []ToggleWithDetail, total int64, err error)
	ListHeadlessAccessPaged(ctx context.Context, in ParamsListHeadlessAccessPaged) (out []ToggleWithDetail, total int64, err error)
	FetchToggleStatByID(ctx context.Context, id string) (out ToggleStat, err error)
	RemoveTogglesByIDs(ctx context.Context, ids []string) (err error)
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
