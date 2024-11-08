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
	"github.com/gorilla/rpc/json"
)

type ConfigRuntime struct {
	Port          int                  `validate:"required"`
	ToggleService moduletoggle.Service `validate:"required"`
}

type Runtime struct {
	Config ConfigRuntime
	Server *ServerStd
}

func NewRuntime(cfg ConfigRuntime) (out *Runtime, err error) {
	if err = validator.Validate(&cfg); err != nil {
		err = fmt.Errorf("bad config: %v", err)
		return
	}

	server := &ServerStd{
		ConfigRuntime: cfg,
		StartedAt:     time.Now(),
	}
	if err = validator.Validate(server); err != nil {
		err = fmt.Errorf("bad server construction: %v", err)
		return
	}

	out = &Runtime{Config: cfg}
	return
}

func (e *Runtime) Run() (err error) {
	s := rpc.NewServer()
	s.RegisterCodec(json.NewCodec(), "application/json")
	s.RegisterService(e.Server, "GoggleRPC")

	r := mux.NewRouter()
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
