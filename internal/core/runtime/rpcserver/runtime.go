package rpcserver

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/avrebarra/goggle/internal/module/moduletoggle"
	"github.com/avrebarra/goggle/utils/validator"
	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
)

type ConfigRuntime struct {
	Version       string               `validate:"required"`
	Port          int                  `validate:"required"`
	ToggleService moduletoggle.Service `validate:"required"`
	StartedAt     time.Time            `validate:"required"`
}

type Runtime struct {
	Config ConfigRuntime
	Server *Handler
}

func NewRuntime(cfg ConfigRuntime) (out *Runtime, err error) {
	cfg.StartedAt = time.Now()
	if err = validator.Validate(&cfg); err != nil {
		err = fmt.Errorf("bad config: %v", err)
		return
	}

	server := &Handler{ConfigRuntime: cfg}

	if err = validator.Validate(server); err != nil {
		err = fmt.Errorf("bad server construction: %v", err)
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
	r.Use(MWRecoverer())
	r.Handle("/", s)

	if err = http.ListenAndServe(fmt.Sprintf(":%d", e.Config.Port), r); err != nil {
		err = fmt.Errorf("error running server: %v", err)
		return
	}

	return
}

func (e *Runtime) Start(ctx context.Context) <-chan bool {
	shutdownChan := make(chan bool)

	go func() {
		<-ctx.Done()
		fmt.Println("shutting down rpc server...")
		close(shutdownChan)
	}()

	go func() {
		log.Printf("starting rpc server in http://localhost:%d\n", e.Config.Port)
		if err := e.Run(); err != nil {
			log.Printf("error running rpc server: %v", err)
		}
	}()

	return shutdownChan
}
