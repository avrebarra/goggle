package utils_test

import (
	"fmt"
	"testing"

	"github.com/avrebarra/goggle/internal/utils"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestExtractStackTrace(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		err := fmt.Errorf("0")
		err = fmt.Errorf("1: %w", err)
		err = fmt.Errorf("2: %w", err)
		err = errors.Wrap(err, "3")
		err = errors.Wrap(err, "4")
		err = errors.Wrap(err, "5")
		err = errors.Wrap(err, "6")
		err = errors.Wrap(err, "7")
		err = errors.Wrap(err, "8")
		err = errors.Wrap(err, "9")

		out, err := utils.ExtractStackTrace(err)
		assert.NoError(t, err)
		assert.NotEmpty(t, out)
	})

	t.Run("has no stacktrace", func(t *testing.T) {
		err := fmt.Errorf("0")
		err = fmt.Errorf("1: %w", err)
		out, err := utils.ExtractStackTrace(err)
		assert.Error(t, err)
		assert.Empty(t, out)
	})

	t.Run("on nil input", func(t *testing.T) {
		out, err := utils.ExtractStackTrace(nil)
		assert.Error(t, err)
		assert.Empty(t, out)
	})
}
