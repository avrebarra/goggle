package httpserver

import (
	"context"
	"time"
)

func (s *Handler) Ping() HandlerFunc {
	type ResponseData struct {
		Version   string    `json:"version"`
		StartedAt time.Time `json:"startedAt"`
		Uptime    string    `json:"uptime"`
	}
	return func(ctx context.Context, rp RequestPack) (err error) {
		out := ResponseData{
			Version:   s.Version,
			StartedAt: s.StartedAt,
			Uptime:    time.Since(s.StartedAt).Round(time.Second).String(),
		}
		return rp.Send(RespSuccess, out)
	}
}
