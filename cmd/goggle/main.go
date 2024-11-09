package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/avrebarra/goggle/internal/module/serviceaccesslog"
	storageaccesslog "github.com/avrebarra/goggle/internal/module/serviceaccesslog/storage"
	storagetoggle "github.com/avrebarra/goggle/internal/module/servicetoggle/storage"

	"github.com/avrebarra/goggle/internal/core/runtime/rpcserver"
	"github.com/avrebarra/goggle/internal/core/runtime/uiserver"
	"github.com/avrebarra/goggle/internal/module/servicetoggle"
	"github.com/avrebarra/goggle/utils/validator"
	"github.com/leaanthony/clir"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	Version = "v0.0.0-unbuilt"
)

type BaseConfig struct {
	PortUI       int    `validate:"required"`
	PortRPC      int    `validate:"required"`
	SQLiteDBPath string `validate:"required"`
}

func main() {
	conf := &BaseConfig{
		PortRPC:      9000,
		PortUI:       9001,
		SQLiteDBPath: "./goggle.db",
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
		rrpc, err := rpcserver.NewRuntime(rpcserver.ConfigRuntime{
			Version:       Version,
			Port:          conf.PortRPC,
			ToggleService: deps.ToggleService,
		})
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

type BaseDeps struct {
	ToggleService servicetoggle.Service
}

func ConstructDeps(conf *BaseConfig) *BaseDeps {
	check := func(err error, name string) {
		if err == nil {
			return
		}
		err = fmt.Errorf("failed to construct dependencies on %s: %v", name, err)
		log.Fatal(err)
	}

	db, err := gorm.Open(sqlite.Open(conf.SQLiteDBPath), &gorm.Config{})
	check(err, "db/sqlite")

	togglestore, err := storagetoggle.NewStorageSQLite(storagetoggle.ConfigStorageSQLite{DB: db})
	check(err, "store/toggle")

	accesslogstore, err := storageaccesslog.NewStorageSQLite(storageaccesslog.ConfigStorageSQLite{DB: db})
	check(err, "store/accesslog")

	accesslogsvc, err := serviceaccesslog.NewService(serviceaccesslog.ServiceConfig{
		AccessLogStore: accesslogstore,
	})
	check(err, "service/accesslog")

	togglesvc, err := servicetoggle.NewService(servicetoggle.ServiceConfig{
		ToggleStore:      togglestore,
		AccessLogService: accesslogsvc,
	})
	check(err, "service/toggle")

	return &BaseDeps{
		ToggleService: togglesvc,
	}
}
