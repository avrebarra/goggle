package do_test

import (
	"context"
	"testing"

	"github.com/avrebarra/goggle/utils/do"
	"github.com/stretchr/testify/assert"
)

func TestParallel(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		ctx := context.Background()
		errs := do.Parallel(ctx, 10, []func() error{
			func() (err error) { return nil },
			func() (err error) { return nil },
			func() (err error) { return nil },
			func() (err error) { return nil },
			func() (err error) { return nil },
		})
		assert.Equal(t, 5, len(errs))

		err := do.JoinErrors(errs)
		assert.NoError(t, err)
	})

	t.Run("on some errored", func(t *testing.T) {
		ctx := context.Background()
		errs := do.Parallel(ctx, 10, []func() error{
			func() (err error) { return nil },
			func() (err error) { return assert.AnError },
			func() (err error) { return nil },
			func() (err error) { return assert.AnError },
			func() (err error) { return nil },
		})
		assert.Equal(t, 5, len(errs))

		err := do.JoinErrors(errs)
		assert.Error(t, err)
	})
}
