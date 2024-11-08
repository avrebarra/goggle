package rpcserver

import (
	"net/http"
	"time"
)

type ServerStd struct {
	ConfigRuntime
	StartedAt time.Time `validate:"required"`
}

func (s *ServerStd) Ping(r *http.Request, in *ReqPing, out *RespPing) error {
	return nil
}

func (s *ServerStd) HealthCheck(r *http.Request, in *ReqHealthCheck, out *RespHealthCheck) error {
	out.Status = "healthy"
	out.StartedAt = s.StartedAt
	out.Uptime = time.Since(s.StartedAt).Round(time.Second).String()
	return nil
}

func (s *ServerStd) Poll(r *http.Request, in *ReqPoll, out *RespPoll) error {
	return nil
}

func (s *ServerStd) GetCompactLogsWithinRange(r *http.Request, in *ReqGetCompactLogsWithinRange, out *RespGetCompactLogsWithinRange) error {
	return nil
}

func (s *ServerStd) GetFullLogsByIDs(r *http.Request, in *ReqGetFullLogsByIDs, out *RespGetFullLogsByIDs) error {
	return nil
}

func (s *ServerStd) FlushFullLog(r *http.Request, in *ReqFlushFullLog, out *RespFlushFullLog) error {
	return nil
}
