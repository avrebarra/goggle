package rpcserver

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/avrebarra/goggle/internal/core/logger"
	"github.com/avrebarra/goggle/internal/utils"
	"github.com/avrebarra/goggle/utils/ctxboard"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

func MWContextSetup() mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = ctxboard.CreateWith(ctx)

			reqctx := RequestContext{StartedAt: time.Now()}
			ctxboard.SetData(ctx, KeyRequestContext, &reqctx)

			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func MWRequestLogger() mux.MiddlewareFunc {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			h.ServeHTTP(w, r.WithContext(ctx))

			reqctx := ctxboard.GetData(ctx, KeyRequestContext).(*RequestContext)
			logdata := logger.OperationLog{
				ID:           reqctx.OpsID,
				Operation:    reqctx.OpsName,
				IngoingData:  reqctx.IngoingData,
				OutgoingData: reqctx.OutgoingData,
				MetaData:     nil,
				ResponseTime: fmt.Sprintf("%dms", time.Since(reqctx.StartedAt).Milliseconds()),
				Error:        nil,
			}
			if reqctx.Error != nil {
				type Error struct {
					Msg        string   `json:"msg"`
					Stacktrace []string `json:"stacktrace"`
				}
				errdata := Error{
					Msg:        reqctx.Error.Error(),
					Stacktrace: []string{},
				}
				if stak, err := utils.ExtractStackTrace(reqctx.Error); err == nil {
					for _, f := range stak {
						errdata.Stacktrace = append(errdata.Stacktrace, fmt.Sprintf("%s in %s", f.FuncName, f.Source))
						if strings.HasPrefix(f.Source, basepath) {
							break
						}
					}
				}
				logdata.Error = errdata
			}

			slog.Info("request finished",
				"type", "opslog",
				"runtime", "rpcserver",
				"opsdata", logdata,
			)
		})
	}
}

func MWRecoverer() mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if recoverable := recover(); recoverable != nil {
					err, isError := recoverable.(error)
					if !isError {
						err = errors.Errorf("%v", recoverable)
					}

					ctx := r.Context()
					reqctx := ctxboard.GetData(ctx, KeyRequestContext).(*RequestContext)
					reqctx.Error = err

					// build output
					var sverr ServerError
					if !errors.As(err, &sverr) {
						sverr = RespErrorPresets[ErrUnexpected]
					}

					resp := &ServerResponse{ID: reqctx.OpsID, Error: sverr}
					w.Header().Set("Content-Type", "application/json; charset=utf-8")
					_ = json.NewEncoder(w).Encode(resp)
				}
			}()
			h.ServeHTTP(w, r)
		})
	}
}
