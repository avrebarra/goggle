package httpserver

import (
	"context"
	"time"

	"github.com/avrebarra/goggle/internal/module/servicetoggle"
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

func (s *Handler) ListToggles() HandlerFunc {
	type ResponseData struct {
		Data  any   `json:"data"`
		Total int64 `json:"total"`
	}
	return func(ctx context.Context, rp RequestPack) (err error) {
		resp, tot, err := s.ToggleService.DoListToggles(ctx, servicetoggle.ParamsDoListToggles{})
		if err != nil {
			err = ErrNotFound.Wrap(err).WithMessage("good days service failure").WithDetail("something failed")
			return
		}

		out := ResponseData{
			Data:  resp,
			Total: tot,
		}
		return rp.Send(RespSuccess, out)
	}
}
