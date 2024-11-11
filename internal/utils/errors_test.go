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
		err := fmt.Errorf("new error")
		err = errors.Wrap(err, "wrapped 1")
		err = errors.Wrapf(err, "wrapped %s", "with format")
		err = errors.Wrapf(err, "wrapped %s", "with format")
		err = errors.Wrapf(err, "wrapped %s", "with format")

		out, err := utils.ExtractStackTrace(err)
		assert.NoError(t, err)
		assert.NotEmpty(t, out)
	})

	t.Run("on nil input", func(t *testing.T) {
		out, err := utils.ExtractStackTrace(nil)
		assert.Error(t, err)
		assert.Empty(t, out)
	})
}
