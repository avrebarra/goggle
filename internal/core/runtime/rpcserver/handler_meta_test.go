package rpcserver_test

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
		req := MakeRPCRequest(t, ServerRequest{
			ID:     "113344",
			Method: "TestRPC.Ping",
			Params: []map[string]any{{}},
		})

		http.HandlerFunc(s.Router.ServeHTTP).ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Code)

		data := res.Body.String()
		getjson := JSONGetter(t, data, true)
		assert.Equal(t, getjson("id").String(), "113344")
		assert.Equal(t, getjson("result.version").String(), "v1.0.0-test")
		assert.Equal(t, getjson("result.startedAt").String(), "2024-11-10T00:00:00Z")
	})
}
