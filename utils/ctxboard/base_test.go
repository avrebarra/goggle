package ctxboard_test

import (
	"context"
	"testing"

	"github.com/avrebarra/goggle/utils/ctxboard"
	"github.com/stretchr/testify/assert"
)

func TestContextBoard_Default(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		out := ctxboard.ContextBoard{}.Default()
		assert.NotNil(t, out)
	})
}

func TestCreateWith(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		ctx := context.Background()
		out := ctxboard.CreateWith(ctx)
		assert.NotNil(t, out)
	})
}

func TestExtractFrom(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		ctx := context.Background()
		ctx = ctxboard.CreateWith(ctx)

		out, ok := ctxboard.ExtractFrom(ctx)
		assert.True(t, ok)
		assert.NotNil(t, out)
	})

	t.Run("on not ok", func(t *testing.T) {
		ctx := context.Background()

		out, ok := ctxboard.ExtractFrom(ctx)
		assert.False(t, ok)
		assert.NotNil(t, out)
	})
}

func TestSetData(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		ctx := context.Background()
		ctx = ctxboard.CreateWith(ctx)

		ctxboard.SetData(ctx, "foo", "bar")
	})
}

func TestGetData(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		ctx := context.Background()
		ctx = ctxboard.CreateWith(ctx)

		ctxboard.SetData(ctx, "foo", "bar")

		out := ctxboard.GetData(ctx, "foo")
		assert.Equal(t, "bar", out)
	})
}
