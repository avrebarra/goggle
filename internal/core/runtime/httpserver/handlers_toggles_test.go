package httpserver_test

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
		req := BakeAPIRequest(t, APIRequest{
			Method: GET,
			Path:   "/list",
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

		seek := JSONSeeker(t, res.Body.String(), true)
		assert.NotEmpty(t, seek("message").String())
		assert.Equal(t, []any{"results", "total"}, seek("data.@keys").Value())
		assert.NotContains(t, seek("data.results.#.id").Value().([]any), "")
	})
}
