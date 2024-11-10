package do_test

import (
	"testing"

	"github.com/avrebarra/goggle/utils/do"
	"github.com/stretchr/testify/assert"
)

func TestJointError_Error(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		je := do.JointError{Errors: []error{assert.AnError, assert.AnError}}
		out := je.Error()
		assert.Equal(t, "2 joint error: assert.AnError general error for testing; assert.AnError general error for testing", out)
	})

	t.Run("on array contains nil", func(t *testing.T) {
		je := do.JointError{Errors: []error{assert.AnError, nil}}
		out := je.Error()
		assert.Equal(t, "1 joint error: assert.AnError general error for testing", out)
	})

	t.Run("on array is empty", func(t *testing.T) {
		je := do.JointError{Errors: []error{}}
		out := je.Error()
		assert.Equal(t, "empty error", out)
	})
}

func TestJoinErrors(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		out := do.JoinErrors([]error{assert.AnError, assert.AnError})
		assert.Error(t, out)
		assert.Regexp(t, "2 joint error: .+", out.Error())
	})

	t.Run("on empty array", func(t *testing.T) {
		out := do.JoinErrors([]error{})
		assert.NoError(t, out)
	})
}
