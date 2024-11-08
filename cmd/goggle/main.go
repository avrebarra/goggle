package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/avrebarra/goggle/internal/core/runtime/rpcserver"
	"github.com/avrebarra/goggle/internal/core/runtime/uiserver"
	"github.com/avrebarra/goggle/utils/validator"
	"github.com/leaanthony/clir"
)

var (
	Version = "v0.0.0-unbuilt"
)

type RuntimeConfig struct {
	PortUI  int `validate:"required"`
	PortRPC int `validate:"required"`
}

func main() {
	conf := &RuntimeConfig{
		PortRPC: 9000,
		PortUI:  9001,
	}

	cli := clir.NewCli("goggle", "goggle manager runner", Version)
	cli.IntFlag("port.rpc", "Port to use for rpc", &conf.PortRPC)
	cli.IntFlag("port.ui", "Port to use for ui", &conf.PortUI)
	cli.Action(func() (err error) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// validate config
		if err := validator.Validate(conf); err != nil {
			log.Fatalf("bad config: %v", err)
		}

		// construct dependencies
		deps := ConstructDeps(conf)
		if err := validator.Validate(deps); err != nil {
			log.Fatalf("bad deps: %v", err)
		}

		// construct runtimes
		rrpc, err := rpcserver.NewRuntime(rpcserver.ConfigRuntime{Port: conf.PortRPC})
		if err != nil {
			err = fmt.Errorf("error creating rpc runtime: %v", err)
			return
		}

		rui, err := uiserver.NewRuntime(uiserver.RuntimeConfig{Port: conf.PortUI})
		if err != nil {
			err = fmt.Errorf("error creating ui runtime: %v", err)
			return
		}

		chWaitRPC := rrpc.Start(ctx)
		chWaitUI := rui.Start(ctx)

		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		<-sigs

		cancel()
		<-chWaitRPC
		<-chWaitUI

		return nil
	})

	if err := cli.Run(); err != nil {
		log.Fatal(err)
	}
}

// ***

type RuntimeDeps struct{}

func ConstructDeps(conf *RuntimeConfig) *RuntimeDeps {
	deps := &RuntimeDeps{}
	return deps
}
