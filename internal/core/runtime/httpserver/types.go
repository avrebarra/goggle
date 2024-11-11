package httpserver

import (
	"time"
)

var (
	KeyRequestContext = "http/request-context"
)

type RequestContext struct {
	OpsID        string
	OpsName      string
	StartedAt    time.Time
	IngoingData  any
	OutgoingData any
	Error        error
}

type Handler struct {
	Config
}
