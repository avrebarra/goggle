package moduletoggle

import "context"

type Service interface {
	Poll(ctx context.Context) (err error)
	GetCompactLogsWithinRange(ctx context.Context, p ParamsGetRange) (out []CompactLog, err error)
	GetFullLogsByIDs(ctx context.Context) (out []FullLog, err error)
	FlushFullLog(ctx context.Context, in FullLog) (err error)
}

type ParamsGetRange struct {
}

type CompactLog struct {
}

type FullLog struct {
}
