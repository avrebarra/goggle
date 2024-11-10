package serviceaccesslog

import (
	"context"

	"github.com/avrebarra/goggle/internal/module/serviceaccesslog/domainaccesslog"
)

type Service interface {
	DoListLogs(ctx context.Context, in ParamsDoListLogs) (out []domainaccesslog.AccessLog, total int64, err error)

	AddAccessLog(ctx context.Context, toggleid string) (err error)
	DeleteAccessLogByToggleID(ctx context.Context, toggleid string) (err error)
}

type ParamsDoListLogs struct {
	Offset    int
	Limit     int
	SortBy    string
	SortOrder string

	FilterToggleIDs []string

	SkipTotal bool
}
