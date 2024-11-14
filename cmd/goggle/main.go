package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/avrebarra/goggle/internal/core/runtime/cronworker"
	"github.com/avrebarra/goggle/internal/core/runtime/httpserver"
	"github.com/avrebarra/goggle/internal/core/runtime/rpcserver"
	"github.com/avrebarra/goggle/internal/core/runtime/uiserver"
	"github.com/avrebarra/goggle/internal/module/clientgithubrepo"
	"github.com/avrebarra/goggle/internal/module/serviceaccesslog"
	"github.com/avrebarra/goggle/internal/module/serviceaccesslog/storeaccesslog"
	"github.com/avrebarra/goggle/internal/module/servicetoggle"
	"github.com/avrebarra/goggle/internal/module/servicetoggle/storetoggle"
	"github.com/avrebarra/goggle/utils/validator"
	"github.com/caarlos0/env/v11"
	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
	"github.com/leaanthony/clir"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	Version = "v0.0.0-unbuilt"
	AppName = "goggle"
	AppDesc = "goggle manager runner"
)

type BaseConfig struct {
	DebugMode      bool   `env:"DEBUG_MODE" yaml:"debug_mode"`
	ConfigFilePath string `env:"CONFIG_PATH" validate:"required"`
	PortUI         int    `yaml:"port_ui" env:"PORT_UI" validate:"required"`
	PortRPC        int    `yaml:"port_rpc" env:"PORT_RPC" validate:"required"`
	PortHTTP       int    `yaml:"port_http" env:"PORT_HTTP" validate:"required"`

	DBMode            string `yaml:"db_mode" env:"DB_MODE" validate:"required,oneof=sqlite postgre"`
	DBSQLitePath      string `yaml:"db_sqlite_path" env:"DB_SQLITE_PATH" validate:"required"`
	DBPostgreHost     string `yaml:"db_postgre_host" env:"DB_POSTGRE_HOST" validate:"required_if=DBMode postgre"`
	DBPostgreUser     string `yaml:"db_postgre_user" env:"DB_POSTGRE_USER" validate:"required_if=DBMode postgre"`
	DBPostgrePassword string `yaml:"db_postgre_password" env:"DB_POSTGRE_PASSWORD" validate:"required_if=DBMode postgre"`
	DBPostgreName     string `yaml:"db_postgre_name" env:"DB_POSTGRE_NAME" validate:"required_if=DBMode postgre"`
	DBPostgrePort     int    `yaml:"db_postgre_port" env:"DB_POSTGRE_PORT" validate:"required_if=DBMode postgre"`

	ClientGitHubTimeout time.Duration `yaml:"client_github_timeout" env:"CLIENT_GITHUB_TIMEOUT" validate:"required"`
	ClientGitHubBaseURL string        `yaml:"client_github_base_url" env:"CLIENT_GITHUB_BASE_URL" validate:"required,endswith=/"`

	RunRPCServer  bool
	RunHTTPServer bool
	RunUIServer   bool
	RunCronWorker bool
}

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	conf := &BaseConfig{
		PortRPC:             9000,
		PortUI:              9001,
		PortHTTP:            9002,
		ConfigFilePath:      "./config.yaml",
		DBMode:              "sqlite",
		DBSQLitePath:        "./local.db",
		ClientGitHubBaseURL: "https://api.github.com/",
		ClientGitHubTimeout: 10 * time.Second,

		RunRPCServer:  false,
		RunHTTPServer: false,
		RunUIServer:   false,
		RunCronWorker: false,
	}

	var exec func() error

	cli := clir.NewCli(AppName, AppDesc, Version)
	cli.BoolFlag("debug", "Set debug mode", &conf.DebugMode)
	cli.StringFlag("config", "Config file path", &conf.ConfigFilePath)
	cli.NewSubCommandInheritFlags("cronworker", "").Action(func() error { conf.RunCronWorker = true; return exec() })
	cli.NewSubCommandInheritFlags("httpserver", "").Action(func() error { conf.RunHTTPServer = true; return exec() })
	cli.NewSubCommandInheritFlags("rpcserver", "").Action(func() error { conf.RunRPCServer = true; return exec() })
	cli.NewSubCommandInheritFlags("uiserver", "").Action(func() error { conf.RunUIServer = true; return exec() })
	cli.AddCommand(clir.NewCommand("help", "").Action(func() error { cli.PrintHelp(); return nil }))
	cli.Action(func() error {
		conf.RunCronWorker = true
		conf.RunRPCServer = true
		conf.RunUIServer = true
		conf.RunHTTPServer = true
		return exec()
	})

	exec = func() (err error) {
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
		rcron, err := cronworker.NewRuntime(cronworker.RuntimeConfig{
			RootContext:  ctx,
			GithubClient: deps.GithubClient,
		})
		ensure(err, "cron runtime")

		rrpc, err := rpcserver.NewRuntime(rpcserver.ConfigRuntime{
			Version:       Version,
			Port:          conf.PortRPC,
			ToggleService: deps.ToggleService,
		})
		ensure(err, "rpc runtime")

		rui, err := uiserver.NewRuntime(uiserver.RuntimeConfig{
			Port:      conf.PortUI,
			DebugMode: conf.DebugMode,
		})
		ensure(err, "ui runtime")

		rhttp, err := httpserver.NewRuntime(httpserver.Config{
			Version:       Version,
			Port:          conf.PortHTTP,
			DebugMode:     conf.DebugMode,
			StartedAt:     time.Now(),
			ToggleService: deps.ToggleService,
		})
		ensure(err, "httpserver runtime")

		// start runtimes
		wchs := []<-chan bool{}
		runtimemap := map[Runtime]bool{
			rcron: conf.RunCronWorker,
			rrpc:  conf.RunRPCServer,
			rui:   conf.RunUIServer,
			rhttp: conf.RunHTTPServer,
		}
		for rt, shouldStart := range runtimemap {
			if !shouldStart {
				continue
			}
			wchs = append(wchs, rt.Start(ctx))
		}

		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		<-sigs

		cancel()
		wg := sync.WaitGroup{}
		for _, wch := range wchs {
			wg.Add(1)
			go func() {
				<-wch
				wg.Done()
			}()
		}
		wg.Wait()

		return nil
	}

	if err := cli.Run(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

// ***

type BaseDeps struct {
	ToggleService servicetoggle.Service   `validate:"required"`
	GithubClient  clientgithubrepo.Client `validate:"required"`
}

func ConstructDeps(conf *BaseConfig) *BaseDeps {
	gormlogger := func() logger.Interface {
		e := logger.Discard
		if conf.DebugMode {
			w := log.New(os.Stdout, "\r\n", log.LstdFlags)
			e = logger.New(w, logger.Config{
				Colorful:             true,
				LogLevel:             logger.Info,
				ParameterizedQueries: false,
			})
		}
		return e
	}()

	var err error
	var togglestore storetoggle.Storage
	var accesslogstore storeaccesslog.Storage
	switch conf.DBMode {
	case "postgre":
		dsnpostgres := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", conf.DBPostgreHost, conf.DBPostgrePort, conf.DBPostgreUser, conf.DBPostgrePassword, conf.DBPostgreName)
		db, err := gorm.Open(postgres.New(postgres.Config{
			DSN:                  dsnpostgres,
			PreferSimpleProtocol: true, // disables implicit prepared statement usage
		}), &gorm.Config{
			Logger: gormlogger,
		})
		ensure(err, "deps db/postgres")

		togglestore, err = storetoggle.NewStoragePostgre(storetoggle.ConfigStoragePostgre{DB: db})
		ensure(err, "deps store/postgres/toggle")

		accesslogstore, err = storeaccesslog.NewStoragePostgre(storeaccesslog.ConfigStoragePostgre{DB: db})
		ensure(err, "deps store/postgres/accesslog")

	default:
		db, err := gorm.Open(sqlite.Open(conf.DBSQLitePath), &gorm.Config{
			SkipDefaultTransaction: true,
			Logger:                 gormlogger,
		})
		ensure(err, "deps db/sqlite")

		togglestore, err = storetoggle.NewStorageSQLite(storetoggle.ConfigStorageSQLite{DB: db})
		ensure(err, "deps store/sqlite/toggle")

		accesslogstore, err = storeaccesslog.NewStorageSQLite(storeaccesslog.ConfigStorageSQLite{DB: db})
		ensure(err, "deps store/sqlite/accesslog")
	}

	httpcli := resty.New()
	httpcli.SetDebug(conf.DebugMode)
	httpcli.SetTimeout(conf.ClientGitHubTimeout)
	clientGithub, err := clientgithubrepo.NewHTTP(clientgithubrepo.ConfigHTTP{
		RESTClient: httpcli,
		BaseURL:    conf.ClientGitHubBaseURL,
	})
	ensure(err, "deps client/github")

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
		GithubClient:  clientGithub,
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

type Runtime interface {
	Start(ctx context.Context) <-chan bool
}

func ensure(err error, name string) {
	if err == nil {
		return
	}
	err = errors.Errorf("ensuring `%s` failed: %v", name, err)
	slog.Error(err.Error())
	os.Exit(1)
}
