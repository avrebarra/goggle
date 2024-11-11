package cronworker

import "context"

type CronHandler struct{ RuntimeConfig }

type CronFunc func(ctx context.Context) (out any, err error)
