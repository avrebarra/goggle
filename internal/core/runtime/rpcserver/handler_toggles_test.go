package rpcserver_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/avrebarra/goggle/internal/module/servicetoggle"
	"github.com/avrebarra/goggle/internal/module/servicetoggle/domaintoggle"
	"github.com/stretchr/testify/assert"
)

func TestHandler_ListToggles(t *testing.T) {
	s := SetupSuite(t)

	t.Run("ok", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := MakeRPCRequest(t, ServerRequest{
			ID:     "113344",
			Method: "TestRPC.ListToggles",
			Params: []map[string]any{{}},
		})

		s.MockToggleService.DoListTogglesFunc = func(ctx context.Context, in servicetoggle.ParamsDoListToggles) ([]domaintoggle.ToggleWithDetail, int64, error) {
			out := []domaintoggle.ToggleWithDetail{}
			for i := 0; i < 3; i++ {
				out = append(out, domaintoggle.ToggleWithDetail{}.Fake())
			}
			return out, int64(len(out)), nil
		}

		http.HandlerFunc(s.Router.ServeHTTP).ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Code)

		data := res.Body.String()
		getjson := JSONGetter(t, data, true)
		assert.Equal(t, getjson("id").String(), "113344")
		assert.Equal(t, len(getjson("result.items").Array()), 3)
	})

	t.Run("on bad params", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := MakeRPCRequest(t, ServerRequest{
			ID:     "113344",
			Method: "TestRPC.ListToggles",
			Params: []map[string]any{{
				"limit": 5000,
			}},
		})

		s.MockToggleService.DoListTogglesFunc = func(ctx context.Context, in servicetoggle.ParamsDoListToggles) ([]domaintoggle.ToggleWithDetail, int64, error) {
			out := []domaintoggle.ToggleWithDetail{}
			for i := 0; i < 3; i++ {
				out = append(out, domaintoggle.ToggleWithDetail{}.Fake())
			}
			return out, int64(len(out)), nil
		}

		http.HandlerFunc(s.Router.ServeHTTP).ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Code)

		data := res.Body.String()
		getjson := JSONGetter(t, data, true)
		assert.Equal(t, getjson("id").String(), "113344")
		assert.Equal(t, getjson("error.code").String(), "unexpected")
		assert.Equal(t, getjson("error.message").String(), "unexpected error: bad request: Limit must be max{100}, actual is 5000")
	})

	t.Run("on service call failure", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := MakeRPCRequest(t, ServerRequest{
			ID:     "113344",
			Method: "TestRPC.ListToggles",
			Params: []map[string]any{{}},
		})

		s.MockToggleService.DoListTogglesFunc = func(ctx context.Context, in servicetoggle.ParamsDoListToggles) ([]domaintoggle.ToggleWithDetail, int64, error) {
			return nil, 0, assert.AnError
		}

		http.HandlerFunc(s.Router.ServeHTTP).ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Code)

		data := res.Body.String()
		getjson := JSONGetter(t, data, true)
		assert.Equal(t, getjson("id").String(), "113344")
		assert.Equal(t, getjson("error.code").String(), "unexpected")
		assert.Equal(t, getjson("error.message").String(), "unexpected error: service failure: assert.AnError general error for testing")
	})
}

func TestHandler_ListStrayToggles(t *testing.T) { t.SkipNow() }
func TestHandler_GetToggle(t *testing.T)        { t.SkipNow() }
func TestHandler_CreateToggle(t *testing.T)     { t.SkipNow() }
func TestHandler_UpdateToggle(t *testing.T)     { t.SkipNow() }
func TestHandler_RemoveToggle(t *testing.T)     { t.SkipNow() }
func TestHandler_StatToggle(t *testing.T)       { t.SkipNow() }
