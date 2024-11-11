package httpserver_test

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/avrebarra/goggle/internal/core/runtime/httpserver"
	"github.com/avrebarra/goggle/utils/validator"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

type TestSuite struct {
	Router            *echo.Echo
	MockToggleService *ToggleServiceMock
}

func SetupSuite(t *testing.T) *TestSuite {
	ts := &TestSuite{}

	slog.SetLogLoggerLevel(slog.LevelError) // disable logs
	gofakeit.Seed(333555444)                // for deterministic tests

	timestamp, err := time.Parse(time.DateOnly, "2024-11-10")
	require.NoError(t, err)

	ts.MockToggleService = &ToggleServiceMock{}

	cfg := httpserver.Config{
		Version:       "v1.0.0-test",
		Port:          3355,
		ToggleService: ts.MockToggleService,
		StartedAt:     timestamp,
	}
	err = validator.Validate(&cfg)
	require.NoError(t, err)

	h := httpserver.Handler{Config: cfg}

	r := echo.New()
	r.Use(middleware.RemoveTrailingSlash())
	r.Use(httpserver.MWCORS())
	r.Use(httpserver.MWContextSetup())
	r.Use(httpserver.MWLogger())
	r.Use(httpserver.MWRecoverer())
	RegisterTestRoutes(t, r, h)
	ts.Router = r

	return ts
}

func RegisterTestRoutes(t *testing.T, r *echo.Echo, h httpserver.Handler) {
	// register routes here:
	r.GET("/", httpserver.Wrap(h.Ping(), "test/ping"))
	r.GET("/list", httpserver.Wrap(h.ListToggles(), "test/list-toggles"))
}

// ***

type HTTPMethod string

var (
	GET    HTTPMethod = "GET"
	POST   HTTPMethod = "POST"
	DELETE HTTPMethod = "DELETE"
	PATCH  HTTPMethod = "PATCH"
	PUT    HTTPMethod = "PUT"
)

type APIRequest struct {
	Method  HTTPMethod        `validate:"required"`
	Headers map[string]string `validate:"-"`
	Path    string            `validate:"required"`
	Data    any               `validate:"-"`
}

func BakeAPIRequest(t *testing.T, in APIRequest) (req *http.Request) {
	data := in.Data

	time.Now().Before(time.Now())

	err := validator.Validate(in)
	if err != nil {
		err = errors.Wrap(err, "cannot bake request due to invalid request details")
	}
	require.NoError(t, err)

	out, err := json.Marshal(data)
	require.NoError(t, err)

	datastr := string(out)
	return httptest.NewRequest(string(in.Method), in.Path, strings.NewReader(datastr))
}

// ***

func JSONSeeker(t *testing.T, base string, serious bool) func(string) gjson.Result {
	if !serious {
		assert.Emptyf(t, base, "the base have this value")
	}
	return func(s string) gjson.Result {
		return gjson.Get(base, s)
	}
}
