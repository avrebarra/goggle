package utils_test

import (
	"testing"
	"unsafe"

	"github.com/avrebarra/goggle/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestUnmarshalToMap(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		out := utils.UnmarshalToMap([]byte(`{"data":123,"values":{"sub1":"ok"}}`))

		assert.NotEmpty(t, out)
		assert.Equal(t, len(out), 2)
	})

	t.Run("err bad data", func(t *testing.T) {
		out := utils.UnmarshalToMap([]byte(`{"data":123,"values`))

		assert.Empty(t, out)
	})
}

func TestRemarshalToMap(t *testing.T) {
	type Data struct {
		Value1 string `json:"value1"`
		Value2 int    `json:"value2"`
	}

	t.Run("ok", func(t *testing.T) {
		data := Data{Value1: "ok", Value2: 123}

		out := utils.RemarshalToMap(data)

		assert.NotEmpty(t, out)
		assert.Equal(t, len(out), 2)
	})

	t.Run("err bad data 1", func(t *testing.T) {
		out := utils.RemarshalToMap("1231")

		assert.Empty(t, out)
	})

	t.Run("err bad data 2", func(t *testing.T) {
		out := utils.RemarshalToMap(unsafe.Pointer(nil))

		assert.Empty(t, out)
	})
}

func TestMorphFrom(t *testing.T) {
	type DataA struct {
		ABC string
		DEF string
		IJK string
	}
	type DataB struct {
		GHI string
		DEF string
		ABC string
	}

	t.Run("ok", func(t *testing.T) {
		source := DataB{GHI: "ghi", DEF: "def", ABC: "abc"}
		targ := DataA{}

		err := utils.MorphFrom(&targ, &source, nil)
		assert.NoError(t, err)
		assert.NotEmpty(t, targ)
		assert.Equal(t, targ.ABC, source.ABC)
		assert.Empty(t, targ.IJK)
	})

	t.Run("ok with supplement", func(t *testing.T) {
		source := DataB{GHI: "ghi", DEF: "def", ABC: "abc"}
		targ := DataA{}

		err := utils.MorphFrom(&targ, &source, &DataA{IJK: "ijkl"})
		assert.NoError(t, err)
		assert.NotEmpty(t, targ)
		assert.Equal(t, targ.ABC, source.ABC)
		assert.Equal(t, targ.IJK, "ijkl")
	})

	t.Run("on bad compose", func(t *testing.T) {
		targ := DataA{}

		err := utils.MorphFrom(&targ, nil, &DataA{IJK: "ijkl"})
		assert.Error(t, err)
	})
}

func TestApplyDefaults(t *testing.T) {
	type DataA struct {
		ABC string
		DEF string
	}

	t.Run("ok", func(t *testing.T) {
		defaults := DataA{ABC: "abc", DEF: "def"}
		data := DataA{}

		utils.ApplyDefaults(&data, &defaults)

		assert.NotEmpty(t, data)
	})
}
