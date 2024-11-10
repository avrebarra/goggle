package ctxsaga_test

import (
	"context"
	"testing"

	"github.com/avrebarra/goggle/utils/ctxboard"
	"github.com/avrebarra/goggle/utils/ctxsaga"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateSaga(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		ctx := ctxboard.CreateWith(context.Background())

		saga := ctxsaga.CreateSaga(ctx)
		assert.NotNil(t, saga)
	})
}

func TestGetSaga(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		ctx := ctxboard.CreateWith(context.Background())
		saga := ctxsaga.CreateSaga(ctx)

		out, ok := ctxsaga.GetSaga(ctx)

		assert.True(t, ok)
		assert.NotNil(t, out)
		assert.Equal(t, saga, out)
	})

	t.Run("on no saga", func(t *testing.T) {
		ctx := ctxboard.CreateWith(context.Background())

		out, ok := ctxsaga.GetSaga(ctx)

		assert.False(t, ok)
		assert.Nil(t, out)
	})
}

func TestSagaCenter_AddRollbackFx(t *testing.T) {
	ctx := ctxboard.CreateWith(context.Background())
	saga := ctxsaga.CreateSaga(ctx)

	t.Run("ok", func(t *testing.T) {
		out, ok := ctxsaga.GetSaga(ctx)
		require.True(t, ok)
		require.Equal(t, out, saga)

		saga.AddRollbackFx(func() error { return nil })
	})
}
func TestSagaCenter_AddCommitFx(t *testing.T) {
	ctx := ctxboard.CreateWith(context.Background())
	saga := ctxsaga.CreateSaga(ctx)

	t.Run("ok", func(t *testing.T) {
		out, ok := ctxsaga.GetSaga(ctx)
		require.True(t, ok)
		require.Equal(t, out, saga)

		saga.AddCommitFx(func() error { return nil })
	})
}

func TestSagaCenter_Commit(t *testing.T) {
	ctx := ctxboard.CreateWith(context.Background())
	saga := ctxsaga.CreateSaga(ctx)

	saga.AddCommitFx(func() error { return nil })
	saga.AddCommitFx(func() error { return nil })

	t.Run("ok", func(t *testing.T) {
		s, ok := ctxsaga.GetSaga(ctx)
		require.True(t, ok)

		err := s.Commit()
		assert.NoError(t, err)
	})

	t.Run("has error", func(t *testing.T) {
		s, ok := ctxsaga.GetSaga(ctx)
		require.True(t, ok)

		s.AddCommitFx(func() error { return assert.AnError })
		s.AddCommitFx(func() error { return assert.AnError })

		err := s.Commit()
		assert.Error(t, err)
	})
}
func TestSagaCenter_Rollback(t *testing.T) {
	ctx := ctxboard.CreateWith(context.Background())
	saga := ctxsaga.CreateSaga(ctx)

	saga.AddRollbackFx(func() error { return nil })
	saga.AddRollbackFx(func() error { return nil })

	t.Run("ok", func(t *testing.T) {
		s, ok := ctxsaga.GetSaga(ctx)
		require.True(t, ok)

		err := s.Rollback()
		assert.NoError(t, err)
	})

	t.Run("has error", func(t *testing.T) {
		s, ok := ctxsaga.GetSaga(ctx)
		require.True(t, ok)

		s.AddRollbackFx(func() error { return assert.AnError })
		s.AddRollbackFx(func() error { return assert.AnError })

		err := s.Rollback()
		assert.Error(t, err)
	})
}
