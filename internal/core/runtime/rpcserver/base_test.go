package rpcserver_test

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/avrebarra/goggle/internal/core/runtime/rpcserver"
	"github.com/avrebarra/goggle/utils/validator"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

type TestSuite struct {
	Router            *mux.Router
	MockToggleService *ToggleServiceMock
}

func SetupSuite(t *testing.T) *TestSuite {
	ts := &TestSuite{}

	slog.SetLogLoggerLevel(slog.LevelError)
	gofakeit.Seed(333555444) // for deterministic tests

	timestamp, err := time.Parse(time.DateOnly, "2024-11-10")
	require.NoError(t, err)

	ts.MockToggleService = &ToggleServiceMock{}

	cfg := rpcserver.ConfigRuntime{
		Version:       "v1.0.0-test",
		Port:          3355,
		ToggleService: ts.MockToggleService,
		StartedAt:     timestamp,
	}
	err = validator.Validate(&cfg)
	require.NoError(t, err)

	s := rpc.NewServer()
	s.RegisterCodec(&rpcserver.Codec{}, "application/json")
	s.RegisterService(&rpcserver.Handler{ConfigRuntime: cfg}, "TestRPC")

	r := mux.NewRouter()
	r.Use(rpcserver.MWContextSetup())
	r.Use(rpcserver.MWRequestLogger())
	r.Use(rpcserver.MWRecoverer())
	r.Handle("/", s)
	ts.Router = r

	return ts
}

func JSONGetter(t *testing.T, base string, serious bool) func(string) gjson.Result {
	if !serious {
		assert.Emptyf(t, base, "the base have this value")
	}
	return func(s string) gjson.Result {
		return gjson.Get(base, s)
	}
}

type ServerRequest struct {
	ID     string `json:"id"`
	Method string `json:"method"`
	Params any    `json:"params"`
}

func MakeRPCRequest(t *testing.T, in ServerRequest) (req *http.Request) {
	out, err := json.Marshal(in)
	require.NoError(t, err)
	datastr := string(out)
	return httptest.NewRequest("POST", "/", strings.NewReader(datastr))
}
