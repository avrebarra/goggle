package rpcserver

import "time"

type ReqPing struct{}
type RespPing struct{}

type ReqHealthCheck struct{}
type RespHealthCheck struct {
	Status    string    `json:"status"`
	StartedAt time.Time `json:"startedAt"`
	Uptime    string    `json:"uptime"`
}

type ReqPoll struct{}
type RespPoll struct{}

type ReqGetCompactLogsWithinRange struct{}
type RespGetCompactLogsWithinRange struct{}

type ReqGetFullLogsByIDs struct{}
type RespGetFullLogsByIDs struct{}

type ReqFlushFullLog struct{}
type RespFlushFullLog struct{}
