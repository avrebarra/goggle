package rpcserver

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/avrebarra/goggle/internal/module/servicetoggle"
	"github.com/avrebarra/goggle/utils/validator"
	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
	"github.com/pkg/errors"
)

// mockable:true
type ToggleService servicetoggle.Service

type ConfigRuntime struct {
	Version       string        `validate:"required"`
	Port          int           `validate:"required"`
	ToggleService ToggleService `validate:"required"`
	StartedAt     time.Time     `validate:"required"`
}

type Runtime struct {
	Config ConfigRuntime
	Server *Handler
}

func NewRuntime(cfg ConfigRuntime) (out *Runtime, err error) {
	cfg.StartedAt = time.Now()
	if err = validator.Validate(&cfg); err != nil {
		err = errors.Errorf("bad config: %v", err)
		return
	}

	server := &Handler{ConfigRuntime: cfg}

	if err = validator.Validate(server); err != nil {
		err = errors.Errorf("bad server construction: %v", err)
		return
	}

	out = &Runtime{Config: cfg}
	return
}

func (e *Runtime) Run() (err error) {
	server := Handler{ConfigRuntime: e.Config}

	s := rpc.NewServer()
	s.RegisterCodec(&Codec{}, "application/json")
	s.RegisterService(&server, "GoggleRPC")

	r := mux.NewRouter()
	r.Use(MWContextSetup())
	r.Use(MWRequestLogger())
	r.Use(MWRecoverer())
	r.Handle("/", s)

	if err = http.ListenAndServe(fmt.Sprintf(":%d", e.Config.Port), r); err != nil {
		err = errors.Errorf("error running server: %v", err)
		return
	}

	return
}

func (e *Runtime) Start(ctx context.Context) <-chan bool {
	shutdownChan := make(chan bool)

	go func() {
		<-ctx.Done()
		slog.Info("shutting down rpcserver...")
		close(shutdownChan)
	}()

	go func() {
		slog.Info(fmt.Sprintf("starting rpcserver in http://localhost:%d", e.Config.Port))
		if err := e.Run(); err != nil {
			err = errors.Wrap(err, "error running rpcserver")
			slog.Error(err.Error())
		}
	}()

	return shutdownChan
}
