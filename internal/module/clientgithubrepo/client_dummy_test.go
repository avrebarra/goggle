package clientgithubrepo_test

import (
	"context"
	"testing"

	"github.com/avrebarra/goggle/internal/module/clientgithubrepo"
	"github.com/stretchr/testify/require"
)

func TestDummy(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		c := clientgithubrepo.Dummy{}

		n, err := c.GetPopularRepoNames(context.Background())
		require.NotEmpty(t, n)
		require.NoError(t, err)
	})
}
