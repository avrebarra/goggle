package cronworker_test

import (
	"context"
	"log/slog"
	"testing"

	"github.com/avrebarra/goggle/internal/core/runtime/cronworker"
	"github.com/avrebarra/goggle/utils/validator"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
)

type TestSuite struct {
	Handler          *cronworker.CronHandler
	MockGitHubClient *GitHubClientMock
}

func SetupSuite(t *testing.T) *TestSuite {
	ts := &TestSuite{}

	slog.SetLogLoggerLevel(slog.LevelError)
	gofakeit.Seed(333555444) // for deterministic tests

	ts.MockGitHubClient = &GitHubClientMock{}

	cfg := cronworker.RuntimeConfig{
		RootContext:  context.Background(),
		GithubClient: ts.MockGitHubClient,
	}
	err := validator.Validate(&cfg)
	require.NoError(t, err)

	ch := cronworker.CronHandler{RuntimeConfig: cfg}
	ts.Handler = &ch

	return ts
}
