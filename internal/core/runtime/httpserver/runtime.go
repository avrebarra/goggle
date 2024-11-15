package httpserver

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/avrebarra/goggle/internal/module/servicetoggle"
	"github.com/avrebarra/goggle/utils/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
)

// mockable:true
type ToggleService servicetoggle.Service

type Runtime struct {
	Config Config
	Server *http.Server
}

type Config struct {
	DebugMode     bool
	Port          int           `validate:"required"`
	Version       string        `validate:"required"`
	StartedAt     time.Time     `validate:"required"`
	ToggleService ToggleService `validate:"required"`
}

func NewRuntime(cfg Config) (out *Runtime, err error) {
	if err = validator.Validate(&cfg); err != nil {
		err = errors.Errorf("bad config")
		return
	}
	out = &Runtime{Config: cfg}
	return
}

func (e *Runtime) Run() (err error) {
	h := Handler{Config: e.Config}

	r := echo.New()

	r.Use(middleware.RemoveTrailingSlash())
	r.Use(MWCORS())
	r.Use(MWContextSetup())
	r.Use(MWLogger())
	r.Use(MWRecoverer())

	r.GET("/", Wrap(h.Ping(), "ping"))
	r.GET("/toggles", Wrap(h.ListToggles(), "list-toggles"))
	r.GET("/toggles/strays", Wrap(h.ListStrayToggles(), "list-stray-toggles"))
	r.GET("/toggles/:id", Wrap(h.GetToggle(), "get-toggle"))
	r.POST("/toggles", Wrap(h.CreateToggle(), "create-toggle"))
	r.PUT("/toggles/:id", Wrap(h.UpdateToggle(), "update-toggle"))
	r.DELETE("/toggles/:id", Wrap(h.RemoveToggle(), "remove-toggle"))
	r.POST("/stat-toggle/:id", Wrap(h.StatToggle(), "stat-toggle"))

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
		slog.Info("shutting down httpserver...")
		e.Server.Shutdown(ctx)
		close(shutdownChan)
	}()

	go func() {
		slog.Info(fmt.Sprintf("starting httpserver in http://localhost:%d", e.Config.Port))
		if err := e.Run(); err != nil {
			err = errors.Wrap(err, "error running httpserver")
			slog.Error(err.Error())
		}
	}()

	return shutdownChan
}
