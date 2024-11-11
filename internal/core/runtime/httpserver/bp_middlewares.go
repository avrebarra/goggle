package httpserver

import (
	"context"
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
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	"github.com/rs/xid"
)

func MWCORS() echo.MiddlewareFunc {
	opts := middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"POST", "GET", "DELETE"},
		AllowHeaders: []string{"Accept", "Content-Type"},
	}

	return middleware.CORSWithConfig(opts)
}

func MWContextSetup() echo.MiddlewareFunc {
	return echo.WrapMiddleware(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = ctxboard.CreateWith(ctx)

			reqctx := RequestContext{StartedAt: time.Now()}
			ctxboard.SetData(ctx, KeyRequestContext, &reqctx)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})
}

func MWLogger() echo.MiddlewareFunc {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	return echo.WrapMiddleware(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)

			ctx := r.Context()
			reqctx := ctxboard.GetData(ctx, KeyRequestContext).(*RequestContext)

			logdata := logger.OperationLog{
				ID:           reqctx.OpsID,
				Operation:    reqctx.OpsName,
				IngoingData:  reqctx.IngoingData,
				OutgoingData: reqctx.OutgoingData,
				MetaData:     nil,
				ResponseTime: fmt.Sprintf("%dms", time.Since(reqctx.StartedAt).Milliseconds()),
				Error:        reqctx.Error,
			}

			err := reqctx.Error
			if err != nil {
				logdata.Error = err
				type Error struct {
					Msg        string   `json:"msg"`
					Stacktrace []string `json:"stacktrace"`
				}
				errdata := Error{
					Msg:        err.Error(),
					Stacktrace: []string{},
				}
				if stak, err := utils.ExtractStackTrace(err); err == nil {
					for _, f := range stak {
						errdata.Stacktrace = append(errdata.Stacktrace, fmt.Sprintf("%s in %s", f.FuncName, f.Source))
						if strings.HasPrefix(f.Source, basepath) {
							break
						}
					}
				}
				logdata.Error = errdata
			}

			slog.Info("http request finished",
				"type", "opslog",
				"runtime", "httpserver",
				"opsdata", logdata,
			)
		})
	})
}

func MWRecoverer() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()
			reqctx := ctxboard.GetData(ctx, KeyRequestContext).(*RequestContext)

			defer func() {
				if recoverable := recover(); recoverable != nil {
					var sverr ServerError
					err, isError := recoverable.(error)
					if !isError {
						err = errors.Errorf("recovered panic: %v", recoverable)
					}
					if !errors.As(err, &sverr) {
						sverr = ServerError{
							Code:  "nocode",
							Cause: "unexpected error",
						}
					}
					expliciterr := errors.Wrap(sverr.error, sverr.Error())
					reqctx.Error = expliciterr

					// build output
					var tpl preset
					switch true {
					case errors.Is(err, ErrNotFound):
						tpl = respmapper[RespNotFound]
					default:
						tpl = respmapper[RespUnexpected]
					}

					status := tpl.Status
					out := tpl.Resp
					out.Error = sverr
					out = out.Normalize()
					reqctx.OutgoingData = out

					c.JSON(status, out)
				}
			}()

			return next(c)
		}
	}
}

// ***

type RequestPack struct {
	echo.Context
}

func (r RequestPack) Bind(t any) (err error) {
	c := r.Context
	ctx := c.Request().Context()
	reqctx := ctxboard.GetData(ctx, KeyRequestContext).(*RequestContext)

	res := r.Context.Bind(t)
	reqctx.IngoingData = res

	return res
}

func (r RequestPack) Send(kind RespKind, data any) (err error) {
	c := r.Context
	ctx := c.Request().Context()
	reqctx := ctxboard.GetData(ctx, KeyRequestContext).(*RequestContext)

	tpl := respmapper[kind]
	status := tpl.Status
	out := tpl.Resp
	out.Data = data
	out = out.Normalize()

	reqctx.OutgoingData = out

	return r.Context.JSON(status, out)
}

type HandlerFunc func(ctx context.Context, pack RequestPack) (err error)

func Wrap(fx HandlerFunc, procname string) func(c echo.Context) error {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		procname = fmt.Sprintf("%s@%s", strings.ToLower(c.Request().Method), procname)

		reqctx := ctxboard.GetData(ctx, KeyRequestContext).(*RequestContext)
		reqctx.OpsName = procname
		reqctx.OpsID = xid.New().String()

		if err := fx(ctx, RequestPack{c}); err != nil {
			panic(err)
		}

		return nil
	}
}

// ***

func HandleNotFound() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		panic(ErrNotFound)
	}
}
