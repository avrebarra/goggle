package rpcserver

import (
	"net/http"
	"time"

	"github.com/avrebarra/goggle/internal/module/servicetoggle"
	"github.com/avrebarra/goggle/internal/module/servicetoggle/domaintoggle"
	"github.com/avrebarra/goggle/internal/utils"
	"github.com/avrebarra/goggle/utils/validator"
	"github.com/guregu/null/v5"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
)

func (s *Handler) ListToggles(r *http.Request, in *ReqListToggles, out *RespListToggles) (err error) {
	if err = validator.Validate(in); err != nil {
		err = errors.Wrap(err, "bad request")
		return
	}

	ctx := r.Context()
	params := servicetoggle.ParamsDoListToggles{}
	_ = utils.MorphFrom(&params, *in, nil)
	data, tot, err := s.ToggleService.DoListToggles(ctx, params)
	if err != nil {
		err = errors.Wrap(err, "service failure")
		return
	}

	resp := &RespListToggles{Items: []ToggleWithDetail{}, Total: tot}
	for _, d := range data {
		val := ToggleWithDetail{}
		utils.MorphFrom(&val, d, nil)
		if d.LastAccessedAt.IsZero() {
			val.LastAccessedAt = null.Time{}
		}
		resp.Items = append(resp.Items, val)
	}
	return copier.Copy(out, resp)
}

func (s *Handler) ListStrayToggles(r *http.Request, in *ReqListStrayToggles, out *RespListStrayToggles) (err error) {
	if err = validator.Validate(in); err != nil {
		err = errors.Wrap(err, "bad request")
		return
	}

	ctx := r.Context()
	params := servicetoggle.ParamsDoListStrayToggles{}
	_ = utils.MorphFrom(&params, *in, nil)
	data, tot, err := s.ToggleService.DoListStrayToggles(ctx, params)
	if err != nil {
		err = errors.Wrap(err, "service failure")
		return
	}

	resp := &RespListStrayToggles{Items: []ToggleStatLog{}, Total: tot}
	for _, d := range data {
		val := ToggleStatLog{}
		utils.MorphFrom(&val, d, nil)
		if d.LastAccessedAt.IsZero() {
			val.LastAccessedAt = null.Time{}
		}
		resp.Items = append(resp.Items, val)
	}

	return copier.Copy(out, resp)
}

func (s *Handler) GetToggle(r *http.Request, in *ReqGetToggle, out *ToggleWithDetail) (err error) {
	if err = validator.Validate(in); err != nil {
		err = errors.Wrap(err, "bad request")
		return
	}

	ctx := r.Context()
	data, err := s.ToggleService.DoGetToggle(ctx, in.ID)
	if errors.Is(err, servicetoggle.ErrNotFound) {
		err = RespErrorPresets[ErrNotFound]
		return
	}
	if err != nil {
		err = errors.Wrap(err, "service failure")
		return
	}

	resp := ToggleWithDetail{}
	utils.MorphFrom(&resp, data, nil)
	if data.LastAccessedAt.IsZero() {
		resp.LastAccessedAt = null.Time{}
	}
	return copier.Copy(out, resp)
}

func (s *Handler) CreateToggle(r *http.Request, in *Toggle, out *Toggle) (err error) {
	in.UpdatedAt = time.Now()
	if err = validator.Validate(in); err != nil {
		err = errors.Wrap(err, "bad request")
		return
	}

	ctx := r.Context()
	data, err := s.ToggleService.DoCreateToggle(ctx, domaintoggle.Toggle(*in))
	if err != nil {
		err = errors.Wrap(err, "service failure")
		return
	}

	resp := Toggle(data)
	return copier.Copy(out, resp)
}

func (s *Handler) UpdateToggle(r *http.Request, in *ReqUpdateToggle, out *Toggle) (err error) {
	in.Data.ID = in.ID
	in.Data.UpdatedAt = time.Now()
	if err = validator.Validate(in); err != nil {
		err = errors.Wrap(err, "bad request")
		return
	}

	ctx := r.Context()
	data, err := s.ToggleService.DoUpdateToggle(ctx, in.ID, domaintoggle.Toggle(in.Data))
	if err != nil {
		err = errors.Wrap(err, "service failure")
		return
	}

	resp := Toggle(data)
	return copier.Copy(out, resp)
}

func (s *Handler) RemoveToggle(r *http.Request, in *ReqRemoveToggle, out *Toggle) (err error) {
	if err = validator.Validate(in); err != nil {
		err = errors.Wrap(err, "bad request")
		return
	}

	ctx := r.Context()
	data, err := s.ToggleService.DoRemoveToggle(ctx, in.ID)
	if errors.Is(err, servicetoggle.ErrNotFound) {
		err = RespErrorPresets[ErrNotFound]
		return
	}
	if err != nil {
		err = errors.Wrap(err, "service failure")
		return
	}

	resp := Toggle(data)
	return copier.Copy(out, resp)
}

func (s *Handler) StatToggle(r *http.Request, in *ReqStatToggle, out *ToggleCompact) (err error) {
	if err = validator.Validate(in); err != nil {
		err = errors.Wrap(err, "bad request")
		return
	}

	ctx := r.Context()
	data, err := s.ToggleService.DoStatToggle(ctx, in.ID)
	if errors.Is(err, servicetoggle.ErrNotFound) {
		err = RespErrorPresets[ErrNotFound]
		return
	}
	if err != nil {
		err = errors.Wrap(err, "service failure")
		return
	}

	resp := ToggleCompact(data)
	return copier.Copy(out, resp)
}
