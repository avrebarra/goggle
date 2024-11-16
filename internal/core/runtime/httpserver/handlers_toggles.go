package httpserver

import (
	"context"
	"time"

	"github.com/avrebarra/goggle/internal/module/servicetoggle"
	"github.com/avrebarra/goggle/internal/module/servicetoggle/domaintoggle"
	"github.com/avrebarra/goggle/internal/utils"
	"github.com/guregu/null/v5"
	"github.com/pkg/errors"
)

// shared dto entities

type Toggle struct {
	ID        string    `json:"id" validate:"gt=0"`
	Status    bool      `json:"status"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type ToggleStatLog struct {
	ID               string    `json:"id"`
	LastAccessedAt   null.Time `json:"lastAccessedAt"`
	AccessFreqWeekly int       `json:"accessFreqWeekly"`
}

type ToggleWithDetail struct {
	ID               string    `json:"id"`
	Status           bool      `json:"status"`
	UpdatedAt        time.Time `json:"updatedAt"`
	LastAccessedAt   null.Time `json:"lastAccessedAt"`
	AccessFreqWeekly int       `json:"accessFreqWeekly"`
}

type ToggleCompact struct {
	ID     string `json:"id"`
	Status bool   `json:"status"`
}

// handlers

func (s *Handler) ListToggles() HandlerFunc {

	type RequestData struct {
		Offset         int      `query:"offset" validate:"min=0"`
		Limit          int      `query:"limit" validate:"min=0,max=100"`
		SortBy         string   `query:"sortBy"`
		SortOrder      string   `query:"sortOrder"`
		FilterIDs      []string `query:"withIds"`
		FilterAccessed bool     `query:"onlyAccessed"`
	}
	type ResponseData struct {
		Results []ToggleWithDetail `json:"results"`
		Total   int64              `json:"total"`
	}
	return func(ctx context.Context, rp RequestPack) (err error) {
		in := RequestData{}
		if err = rp.Bind(&in); err != nil {
			return
		}

		params := servicetoggle.ParamsDoListToggles{}
		_ = utils.MorphFrom(&params, in, nil)
		resp, tot, err := s.ToggleService.DoListToggles(ctx, params)
		if err != nil {
			err = errors.Wrap(err, "service failure")
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
	type RequestData struct {
		Offset    int    `query:"offset"`
		Limit     int    `query:"limit"`
		SortBy    string `query:"sortBy"`
		SortOrder string `query:"sortOrder"`
	}
	type ResponseData struct {
		Items []ToggleStatLog `json:"items"`
		Total int64           `json:"total"`
	}
	return func(ctx context.Context, rp RequestPack) (err error) {
		in := RequestData{}
		if err = rp.Bind(&in); err != nil {
			return
		}

		params := servicetoggle.ParamsDoListStrayToggles{}
		_ = utils.MorphFrom(&params, in, nil)
		data, tot, err := s.ToggleService.DoListStrayToggles(ctx, params)
		if err != nil {
			err = errors.Wrap(err, "service failure")
			return
		}

		out := ResponseData{Items: []ToggleStatLog{}, Total: tot}
		for _, d := range data {
			val := ToggleStatLog{}
			utils.MorphFrom(&val, d, nil)
			if d.LastAccessedAt.IsZero() {
				val.LastAccessedAt = null.Time{}
			}
			out.Items = append(out.Items, val)
		}

		return rp.Send(RespSuccess, out)
	}
}

func (s *Handler) GetToggle() HandlerFunc {
	type RequestData struct {
		ID string `param:"id"`
	}
	type ResponseData ToggleWithDetail
	return func(ctx context.Context, rp RequestPack) (err error) {
		in := RequestData{}
		if err = rp.Bind(&in); err != nil {
			return
		}

		data, err := s.ToggleService.DoGetToggle(ctx, in.ID)
		if errors.Is(err, servicetoggle.ErrNotFound) {
			err = ErrNotFound
			return
		}
		if err != nil {
			err = errors.Wrap(err, "service failure")
			return
		}

		out := ResponseData{}
		utils.MorphFrom(&out, data, nil)
		if data.LastAccessedAt.IsZero() {
			out.LastAccessedAt = null.Time{}
		}
		return rp.Send(RespSuccess, out)
	}
}

func (s *Handler) CreateToggle() HandlerFunc {
	type RequestData Toggle
	type ResponseData Toggle
	return func(ctx context.Context, rp RequestPack) (err error) {
		in := RequestData{}
		if err = rp.Bind(&in); err != nil {
			return
		}

		data, err := s.ToggleService.DoCreateToggle(ctx, domaintoggle.Toggle(in))
		if err != nil {
			return
		}

		out := ResponseData(data)
		return rp.Send(RespSuccess, out)
	}
}

func (s *Handler) UpdateToggle() HandlerFunc {
	type RequestData struct {
		Toggle
		ID string `param:"id"`
	}
	type ResponseData Toggle
	return func(ctx context.Context, rp RequestPack) (err error) {
		in := RequestData{}
		if err = rp.Bind(&in); err != nil {
			return
		}

		data, err := s.ToggleService.DoUpdateToggle(ctx, in.ID, domaintoggle.Toggle(in.Toggle))
		if err != nil {
			err = errors.Wrap(err, "service failure")
			return
		}

		out := ResponseData(data)
		return rp.Send(RespSuccess, out)
	}
}

func (s *Handler) RemoveToggle() HandlerFunc {
	type RequestData struct {
		ID string `param:"id"`
	}
	type ResponseData Toggle
	return func(ctx context.Context, rp RequestPack) (err error) {
		in := RequestData{}
		if err = rp.Bind(&in); err != nil {
			return
		}

		data, err := s.ToggleService.DoRemoveToggle(ctx, in.ID)
		if errors.Is(err, servicetoggle.ErrNotFound) {
			err = ErrNotFound
			return
		}
		if err != nil {
			err = errors.Wrap(err, "service failure")
			return
		}

		out := ResponseData(data)
		return rp.Send(RespSuccess, out)
	}
}

func (s *Handler) StatToggle() HandlerFunc {
	type RequestData struct {
		ID string `param:"id"`
	}
	type ResponseData ToggleCompact
	return func(ctx context.Context, rp RequestPack) (err error) {
		in := RequestData{}
		if err = rp.Bind(&in); err != nil {
			return
		}

		data, err := s.ToggleService.DoStatToggle(ctx, in.ID)
		if errors.Is(err, servicetoggle.ErrNotFound) {
			err = ErrNotFound
			return
		}
		if err != nil {
			err = errors.Wrap(err, "service failure")
			return
		}

		out := ResponseData(data)
		return rp.Send(RespSuccess, out)
	}
}
