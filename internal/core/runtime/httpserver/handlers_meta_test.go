package httpserver_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandler_Ping(t *testing.T) {
	s := SetupSuite(t)

	t.Run("ok", func(t *testing.T) {
		res := httptest.NewRecorder()
		req := BakeAPIRequest(t, APIRequest{
			Method: GET,
			Path:   "/",
		})

		http.HandlerFunc(s.Router.ServeHTTP).ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Code)

		seek := JSONSeeker(t, res.Body.String(), true)
		assert.NotEmpty(t, seek("message").String())
		assert.Equal(t, []any{"version", "startedAt", "uptime"}, seek("data.@keys").Value())
		assert.Equal(t, "v1.0.0-test", seek("data.version").String())
	})
}
