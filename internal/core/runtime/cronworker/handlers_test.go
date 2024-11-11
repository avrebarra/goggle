package cronworker_test

import (
	"context"
	"testing"

	"github.com/avrebarra/goggle/internal/module/clientgithubrepo"
	"github.com/avrebarra/goggle/utils/ctxboard"
	"github.com/stretchr/testify/assert"
)

func TestCronHandler_ExecGitHubPing(t *testing.T) {
	s := SetupSuite(t)
	ctx := ctxboard.CreateWith(s.Handler.RootContext)

	t.Run("ok", func(t *testing.T) {
		out, err := s.Handler.ExecGitHubPing(ctx)
		assert.NoError(t, err)
		assert.NotNil(t, out)
	})

	t.Run("on remote call failure", func(t *testing.T) {
		s.MockGitHubClient.GetTopRepoDetailsFunc = func(ctx context.Context) ([]clientgithubrepo.RepoDetail, error) {
			return nil, assert.AnError
		}

		out, err := s.Handler.ExecGitHubPing(ctx)
		assert.Error(t, err)
		assert.Nil(t, out)
	})
}
