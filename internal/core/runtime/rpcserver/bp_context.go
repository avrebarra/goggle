package rpcserver

import "time"

var (
	KeyRequestContext = "rpc/request-context"
)

type RequestContext struct {
	OpsID        string
	OpsName      string
	StartedAt    time.Time
	IngoingData  any
	OutgoingData any
	Error        error
}
