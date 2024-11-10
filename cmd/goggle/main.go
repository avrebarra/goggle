package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/avrebarra/goggle/internal/core/runtime/rpcserver"
	"github.com/avrebarra/goggle/internal/core/runtime/uiserver"
	"github.com/avrebarra/goggle/internal/module/serviceaccesslog"
	"github.com/avrebarra/goggle/internal/module/serviceaccesslog/storeaccesslog"
	"github.com/avrebarra/goggle/internal/module/servicetoggle"
	"github.com/avrebarra/goggle/internal/module/servicetoggle/storetoggle"
	"github.com/avrebarra/goggle/utils/validator"
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/leaanthony/clir"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	Version = "v0.0.0-unbuilt"
	AppName = "goggle"
	AppDesc = "goggle manager runner"
)

type BaseConfig struct {
	ConfigFilePath string `env:"CONFIG_PATH" validate:"required"`
	PortUI         int    `yaml:"port_ui" env:"PORT_UI" validate:"required"`
	PortRPC        int    `yaml:"port_rpc" env:"PORT_RPC" validate:"required"`
	SQLiteDBPath   string `yaml:"sqlite_db_path" env:"SQLITE_DB_PATH" validate:"required"`
}

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	conf := &BaseConfig{
		PortRPC:        9000,
		PortUI:         9001,
		ConfigFilePath: "./config.yaml",
		SQLiteDBPath:   "./local.db",
	}

	cli := clir.NewCli(AppName, AppDesc, Version)
	cli.StringFlag("config", "Config file path", &conf.ConfigFilePath)
	cli.Action(func() (err error) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// construct config
		conf = ConstructConfig(conf)
		err = validator.Validate(conf)
		ensure(err, "config")

		// construct dependencies
		deps := ConstructDeps(conf)
		err = validator.Validate(deps)
		ensure(err, "dependency")

		// construct runtimes
		rrpc, err := rpcserver.NewRuntime(rpcserver.ConfigRuntime{
			Version:       Version,
			Port:          conf.PortRPC,
			ToggleService: deps.ToggleService,
		})
		ensure(err, "rpc runtime")

		rui, err := uiserver.NewRuntime(uiserver.RuntimeConfig{Port: conf.PortUI})
		ensure(err, "ui runtime")

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
	db, err := gorm.Open(sqlite.Open(conf.SQLiteDBPath), &gorm.Config{})
	ensure(err, "deps db/sqlite")

	togglestore, err := storetoggle.NewStorageSQLite(storetoggle.ConfigStorageSQLite{DB: db})
	ensure(err, "deps store/toggle")

	accesslogstore, err := storeaccesslog.NewStorageSQLite(storeaccesslog.ConfigStorageSQLite{DB: db})
	ensure(err, "deps store/accesslog")

	accesslogsvc, err := serviceaccesslog.NewService(serviceaccesslog.ServiceConfig{
		AccessLogStore: accesslogstore,
	})
	ensure(err, "deps service/accesslog")

	togglesvc, err := servicetoggle.NewService(servicetoggle.ServiceConfig{
		ToggleStore:      togglestore,
		AccessLogService: accesslogsvc,
	})
	ensure(err, "deps service/toggle")

	return &BaseDeps{
		ToggleService: togglesvc,
	}
}

func ConstructConfig(conf *BaseConfig) *BaseConfig {
	switch true {
	case godotenv.Load() == nil:
	case godotenv.Load("local.env") == nil:
	}

	err := env.Parse(conf)
	ensure(err, "env loading")

	_, err = os.Stat(conf.ConfigFilePath)
	if err == nil {
		cfgdata, err := os.ReadFile(conf.ConfigFilePath)
		ensure(err, "configfile reading")

		err = yaml.Unmarshal(cfgdata, &conf)
		ensure(err, "configfile unmarshaling")
	}

	return conf
}

// ***

func ensure(err error, name string) {
	if err == nil {
		return
	}
	err = errors.Errorf("failed to ensure %s: failure on %v", name, err)
	slog.Error(err.Error())
	os.Exit(1)
}
