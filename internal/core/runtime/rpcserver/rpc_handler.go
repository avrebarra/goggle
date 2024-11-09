package rpcserver

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	domaintoggle "github.com/avrebarra/goggle/internal/module/servicetoggle/domain"

	"github.com/avrebarra/goggle/internal/module/servicetoggle"
	"github.com/avrebarra/goggle/internal/utils"
	"github.com/avrebarra/goggle/utils/validator"
	"github.com/guregu/null/v5"
	"github.com/jinzhu/copier"
)

type Handler struct {
	ConfigRuntime
}

// func (s *Server) Sample(r *http.Request, in *ReqPing, out *RespPing) (err error) { return nil }

func (s *Handler) ListToggles(r *http.Request, in *ReqListToggles, out *RespListToggles) (err error) {
	if err = validator.Validate(in); err != nil {
		err = fmt.Errorf("bad request: %w", err)
		return
	}

	ctx := r.Context()
	params := servicetoggle.ParamsDoListToggles{}
	_ = utils.Translate(&params, *in, nil)
	data, tot, err := s.ToggleService.DoListToggles(ctx, params)
	if err != nil {
		err = fmt.Errorf("service failure: %w", err)
		return
	}

	resp := &RespListToggles{Items: []ToggleWithDetail{}, Total: tot}
	for _, d := range data {
		val := ToggleWithDetail{}
		utils.Translate(&val, d, nil)
		if d.LastAccessedAt.IsZero() {
			val.LastAccessedAt = null.Time{}
		}
		resp.Items = append(resp.Items, val)
	}
	return copier.Copy(out, resp)
}

func (s *Handler) ListStrayToggles(r *http.Request, in *ReqListStrayToggles, out *RespListStrayToggles) (err error) {
	if err = validator.Validate(in); err != nil {
		err = fmt.Errorf("bad request: %w", err)
		return
	}

	ctx := r.Context()
	params := servicetoggle.ParamsDoListStrayToggles{}
	_ = utils.Translate(&params, *in, nil)
	data, tot, err := s.ToggleService.DoListStrayToggles(ctx, params)
	if err != nil {
		err = fmt.Errorf("service failure: %w", err)
		return
	}

	resp := &RespListStrayToggles{Items: []ToggleStatLog{}, Total: tot}
	for _, d := range data {
		val := ToggleStatLog{}
		utils.Translate(&val, d, nil)
		if d.LastAccessedAt.IsZero() {
			val.LastAccessedAt = null.Time{}
		}
		resp.Items = append(resp.Items, val)
	}

	return copier.Copy(out, resp)
}

func (s *Handler) GetToggle(r *http.Request, in *ReqGetToggle, out *Toggle) (err error) {
	if err = validator.Validate(in); err != nil {
		err = fmt.Errorf("bad request: %w", err)
		return
	}

	ctx := r.Context()
	data, err := s.ToggleService.DoGetToggle(ctx, in.ID)
	if errors.Is(err, servicetoggle.ErrNotFound) {
		err = RespErrorPresets[ErrDataNotFound]
		return
	}
	if err != nil {
		err = fmt.Errorf("service failure: %w", err)
		return
	}

	resp := ToggleWithDetail{}
	utils.Translate(&resp, data, nil)
	if data.LastAccessedAt.IsZero() {
		resp.LastAccessedAt = null.Time{}
	}
	return copier.Copy(out, resp)
}

func (s *Handler) CreateToggle(r *http.Request, in *Toggle, out *Toggle) (err error) {
	if err = validator.Validate(in); err != nil {
		err = fmt.Errorf("bad request: %w", err)
		return
	}

	ctx := r.Context()
	data, err := s.ToggleService.DoCreateToggle(ctx, domaintoggle.Toggle(*in))
	if err != nil {
		err = fmt.Errorf("service failure: %w", err)
		return
	}

	resp := Toggle(data)
	return copier.Copy(out, resp)
}

func (s *Handler) UpdateToggle(r *http.Request, in *ReqUpdateToggle, out *Toggle) (err error) {
	if err = validator.Validate(in); err != nil {
		err = fmt.Errorf("bad request: %w", err)
		return
	}

	ctx := r.Context()
	data, err := s.ToggleService.DoUpdateToggle(ctx, in.ID)
	if err != nil {
		err = fmt.Errorf("service failure: %w", err)
		return
	}

	resp := Toggle(data)
	return copier.Copy(out, resp)
}

func (s *Handler) RemoveToggle(r *http.Request, in *ReqRemoveToggle, out *Toggle) (err error) {
	if err = validator.Validate(in); err != nil {
		err = fmt.Errorf("bad request: %w", err)
		return
	}

	ctx := r.Context()
	data, err := s.ToggleService.DoRemoveToggle(ctx, in.ID)
	if errors.Is(err, servicetoggle.ErrNotFound) {
		err = RespErrorPresets[ErrDataNotFound]
		return
	}
	if err != nil {
		err = fmt.Errorf("service failure: %w", err)
		return
	}

	resp := Toggle(data)
	return copier.Copy(out, resp)
}

func (s *Handler) StatToggle(r *http.Request, in *ReqStatToggle, out *ToggleCompact) (err error) {
	if err = validator.Validate(in); err != nil {
		err = fmt.Errorf("bad request: %w", err)
		return
	}

	ctx := r.Context()
	data, err := s.ToggleService.DoStatToggle(ctx, in.ID)
	if errors.Is(err, servicetoggle.ErrNotFound) {
		err = RespErrorPresets[ErrDataNotFound]
		return
	}
	if err != nil {
		err = fmt.Errorf("service failure: %w", err)
		return
	}

	resp := ToggleCompact(data)
	return copier.Copy(out, resp)
}

func (s *Handler) Ping(r *http.Request, in *ReqPing, out *RespPing) (err error) {
	resp := RespPing{
		Version:   s.Version,
		StartedAt: s.StartedAt,
		Uptime:    time.Since(s.StartedAt).Round(time.Second).String(),
	}
	return copier.Copy(out, resp)
}
