package uiserver

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/avrebarra/goggle/ui"
	"github.com/avrebarra/goggle/utils/validator"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type Runtime struct {
	Config RuntimeConfig
	Server *http.Server
}

type RuntimeConfig struct {
	DebugMode bool
	Port      int `validate:"required"`
}

func NewRuntime(cfg RuntimeConfig) (out *Runtime, err error) {
	if err = validator.Validate(&cfg); err != nil {
		err = errors.Errorf("bad config")
		return
	}

	out = &Runtime{Config: cfg}
	return
}

func (e *Runtime) Run() (err error) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	if e.Config.DebugMode {
		gin.SetMode(gin.DebugMode)
		r = gin.Default()
	}

	ui.AddRoutes(r)
	e.Server = &http.Server{
		Addr:    fmt.Sprintf(":%d", e.Config.Port),
		Handler: http.HandlerFunc(r.ServeHTTP),
	}

	err = e.Server.ListenAndServe()
	if err == http.ErrServerClosed {
		err = nil
	}
	if err != nil {
		err = errors.Errorf("error running server: %v", err)
		return
	}

	return
}

func (e *Runtime) Start(ctx context.Context) <-chan bool {
	shutdownChan := make(chan bool)

	go func() {
		<-ctx.Done()
		slog.Info("shutting down uiserver...")
		e.Server.Shutdown(ctx)
		close(shutdownChan)
	}()

	go func() {
		slog.Info(fmt.Sprintf("starting uiserver in http://localhost:%d", e.Config.Port))
		if err := e.Run(); err != nil {
			err = errors.Wrap(err, "error running uiserver")
			slog.Error(err.Error())
		}
	}()

	return shutdownChan
}
