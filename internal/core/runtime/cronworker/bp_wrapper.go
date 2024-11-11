package cronworker

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/avrebarra/goggle/internal/core/logger"
	"github.com/avrebarra/goggle/internal/utils"
	"github.com/avrebarra/goggle/utils/ctxboard"
	"github.com/rs/xid"
)

func Wrap(e *Runtime, name string, exec CronFunc, crondef string) (string, func()) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	return crondef, func() {
		ctx := e.Config.RootContext
		timestamp := time.Now()
		procname := strings.ToLower(name)

		ctx = ctxboard.CreateWith(ctx)

		logdata := logger.OperationLog{
			ID:           xid.New().String(),
			Operation:    procname,
			IngoingData:  nil,
			OutgoingData: nil,
			MetaData:     nil,
			ResponseTime: "",
			Error:        nil,
		}

		var err error
		var out any

		defer func() {
			logdata.ResponseTime = fmt.Sprintf("%dms", time.Since(timestamp).Milliseconds())
			logdata.OutgoingData = out

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

			slog.Info("cron job finished",
				"type", "opslog",
				"runtime", "cronworker",
				"opsdata", logdata,
			)
		}()

		out, err = exec(ctx)
	}
}
