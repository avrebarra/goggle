package validator_test

import (
	"testing"

	"github.com/avrebarra/goggle/utils/validator"
	origvalidator "github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

type Data struct {
	Value string `validate:"required,len=4" alias:"isi"`
}

var orig = origvalidator.New()
var sampledata = Data{Value: "testx"}

func TestGlobal(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		out := validator.Global()
		assert.NotNil(t, out)
	})
}

func TestSetGlobal(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		v := validator.Global()
		validator.SetGlobal(v)
	})
}

func TestValidate(t *testing.T) {
	t.SkipNow()
	t.Run("ok", func(t *testing.T) {
		data := sampledata
		err := validator.Validate(data)
		assert.Error(t, err)
	})
}

func TestNew(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		out := validator.New(orig)
		assert.NotNil(t, out)
	})
}

func TestValidator_Validate(t *testing.T) {
	t.Run("ok err", func(t *testing.T) {
		err := validator.Validate(sampledata)
		assert.Error(t, err)
	})

	t.Run("ok no err", func(t *testing.T) {
		err := validator.Validate(Data{Value: "test"})
		assert.NoError(t, err)
	})
}

func TestValidationError_GetRootError(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		err := validator.Validate(sampledata)
		ve := err.(validator.ValidationError)

		out := ve.GetRootError()
		assert.NotNil(t, out)
	})
}

func TestValidationError_Error(t *testing.T) {
	err := validator.Validate(sampledata)
	ve := err.(validator.ValidationError)

	out := ve.Error()
	assert.Equal(t, out, "isi must be len{4}, actual is testx")
	assert.NotNil(t, out)
}

func TestValidationError_Unwrap(t *testing.T) {
	err := validator.Validate(sampledata)
	ve := err.(validator.ValidationError)

	out := ve.Unwrap()
	assert.NotNil(t, out)
}

func TestGetErrorsData(t *testing.T) {
	type DataStructure struct {
		Name string `validate:"required" alias:"nama"`
		Age  int    `validate:"required,gte=0" alias:"usia"`
	}

	t.Run("ok", func(t *testing.T) {
		data := DataStructure{Name: "", Age: 0}

		err := validator.Validate(data)
		verr, ok := validator.ExtractValidationErrors(err)

		assert.True(t, ok)
		assert.NotNil(t, verr)
		assert.NotEmpty(t, verr)
	})

	t.Run("on non validation error", func(t *testing.T) {
		err := assert.AnError
		verr, ok := validator.ExtractValidationErrors(err)

		assert.False(t, ok)
		assert.NotNil(t, verr)
		assert.NotNil(t, verr)
		assert.Empty(t, verr)
	})
}
