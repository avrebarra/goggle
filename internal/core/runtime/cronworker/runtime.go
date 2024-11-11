package cronworker

import (
	"context"
	"log/slog"

	"github.com/avrebarra/goggle/internal/module/clientgithubrepo"
	"github.com/avrebarra/goggle/utils/validator"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
)

// mockable:true
type GitHubClient clientgithubrepo.Client

type Runtime struct {
	Config RuntimeConfig
	Cron   *cron.Cron
}

type RuntimeConfig struct {
	RootContext  context.Context `validate:"required"`
	GithubClient GitHubClient    `validate:"required"`
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
	c := cron.New()
	e.Cron = c

	// setup handlers
	ch := CronHandler{RuntimeConfig: e.Config}
	c.AddFunc(Wrap(e, "cron/ping-github", ch.ExecGitHubPing, "@every 300s"))

	// start engine
	c.Start()

	return
}

func (e *Runtime) Start(ctx context.Context) <-chan bool {
	shutdownChan := make(chan bool)

	go func() {
		<-ctx.Done()
		slog.Info("shutting down cronworker...")
		ctxStop := e.Cron.Stop()
		<-ctxStop.Done()
		close(shutdownChan)
	}()

	go func() {
		slog.Info("starting cronworker")
		if err := e.Run(); err != nil {
			err = errors.Wrap(err, "error running cronworker")
			slog.Error(err.Error())
		}
	}()

	return shutdownChan
}
