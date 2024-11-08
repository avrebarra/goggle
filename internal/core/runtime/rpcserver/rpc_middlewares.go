package rpcserver

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func MWRecoverer() mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				r := recover()
				if r != nil {
					err := fmt.Errorf("panic: %v", r)
					errid := json.RawMessage([]byte(`-1`))
					resp := &ServerResponse{ID: &errid, Error: RespErrorPresets[ErrUnexpected].WithMessage(err.Error())}
					w.Header().Set("Content-Type", "application/json; charset=utf-8")
					_ = json.NewEncoder(w).Encode(resp)
				}
			}()
			h.ServeHTTP(w, r)
		})
	}
}
