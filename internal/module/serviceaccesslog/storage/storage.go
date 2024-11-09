package storage

import (
	"context"
	"errors"

	domainaccesslog "github.com/avrebarra/goggle/internal/module/serviceaccesslog/domain"
)

var (
	ErrStoreNotFound = errors.New("not found")
)

type Storage interface {
	FetchPaged(ctx context.Context, in ParamsFetchPaged) (out []domainaccesslog.AccessLog, total int64, err error)
	CreateLog(ctx context.Context, in domainaccesslog.AccessLog) (out domainaccesslog.AccessLog, err error)
	DeleteAllByToggleIDs(ctx context.Context, toggleids []string) (err error)
}

type ParamsFetchPaged struct {
	Offset    int    `validate:"min=0"`
	Limit     int    `validate:"min=0,max=100"`
	SortBy    string `validate:"oneof=id name"`
	SortOrder string `validate:"oneof=asc desc"`

	FilterToggleIDs []string `validate:"-"`

	SkipTotal bool `validate:"-"`
}
