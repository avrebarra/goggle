package httpserver

import (
	"context"
	"errors"
	"time"

	"github.com/avrebarra/goggle/internal/module/servicetoggle"
	"github.com/avrebarra/goggle/internal/utils"
	"github.com/guregu/null/v5"
)

func (s *Handler) ListToggles() HandlerFunc {
	type ToggleWithDetail struct {
		ID               string    `json:"id"`
		Status           bool      `json:"status"`
		UpdatedAt        time.Time `json:"updatedAt"`
		LastAccessedAt   null.Time `json:"lastAccessedAt"`
		AccessFreqWeekly int       `json:"accessFreqWeekly"`
	}
	type ResponseData struct {
		Results []ToggleWithDetail `json:"results"`
		Total   int64              `json:"total"`
	}
	return func(ctx context.Context, rp RequestPack) (err error) {
		resp, tot, err := s.ToggleService.DoListToggles(ctx, servicetoggle.ParamsDoListToggles{})
		if err != nil {
			return
		}

		results := []ToggleWithDetail{}
		for _, d := range resp {
			val := ToggleWithDetail{}
			utils.MorphFrom(&val, d, nil)
			if d.LastAccessedAt.IsZero() {
				val.LastAccessedAt = null.Time{}
			}
			results = append(results, val)
		}

		out := ResponseData{Results: results, Total: tot}
		return rp.Send(RespSuccess, out)
	}
}

func (s *Handler) ListStrayToggles() HandlerFunc {
	return func(ctx context.Context, rp RequestPack) (err error) {
		return errors.New("not implemented")
	}
}

func (s *Handler) GetToggle() HandlerFunc {
	return func(ctx context.Context, rp RequestPack) (err error) {
		return errors.New("not implemented")
	}
}

func (s *Handler) CreateToggle() HandlerFunc {
	return func(ctx context.Context, rp RequestPack) (err error) {
		return errors.New("not implemented")
	}
}

func (s *Handler) UpdateToggle() HandlerFunc {
	return func(ctx context.Context, rp RequestPack) (err error) {
		return errors.New("not implemented")
	}
}

func (s *Handler) RemoveToggle() HandlerFunc {
	return func(ctx context.Context, rp RequestPack) (err error) {
		return errors.New("not implemented")
	}
}

func (s *Handler) StatToggle() HandlerFunc {
	return func(ctx context.Context, rp RequestPack) (err error) {
		return errors.New("not implemented")
	}
}
