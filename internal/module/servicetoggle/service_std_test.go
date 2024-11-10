package servicetoggle_test

import (
	"context"
	"testing"

	"github.com/avrebarra/goggle/internal/module/servicetoggle"
	"github.com/avrebarra/goggle/internal/module/servicetoggle/domaintoggle"
	"github.com/avrebarra/goggle/internal/module/servicetoggle/storetoggle"
	"github.com/avrebarra/goggle/utils/ctxboard"
	"github.com/avrebarra/goggle/utils/ctxsaga"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewService(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		out, err := servicetoggle.NewService(servicetoggle.ServiceConfig{
			ToggleStore:      &StorageToggleMock{},
			AccessLogService: &ServiceAccessLogMock{},
		})
		assert.NoError(t, err)
		assert.NotNil(t, out)
	})

	t.Run("err on bad deps", func(t *testing.T) {
		_, err := servicetoggle.NewService(servicetoggle.ServiceConfig{})
		assert.Error(t, err)
	})
}

// ***

type TestSuite struct {
	Context              context.Context
	MockToggleStore      *StorageToggleMock
	MockAccessLogService *ServiceAccessLogMock
	ServiceStd           *servicetoggle.ServiceStd
}

func SetupSuite(t *testing.T) *TestSuite {
	s := &TestSuite{}

	s.Context = ctxboard.CreateWith(context.Background())

	s.MockToggleStore = &StorageToggleMock{}
	s.MockAccessLogService = &ServiceAccessLogMock{}

	svc, err := servicetoggle.NewService(servicetoggle.ServiceConfig{
		ToggleStore:      s.MockToggleStore,
		AccessLogService: s.MockAccessLogService,
	})
	assert.NoError(t, err)
	s.ServiceStd = svc

	return s
}

func TestServiceStd_DoListToggles(t *testing.T) {
	setupTest := func() *TestSuite {
		s := SetupSuite(t)

		s.MockToggleStore.FetchPagedFunc = func(ctx context.Context, in storetoggle.ParamsFetchPaged) ([]domaintoggle.ToggleWithDetail, int64, error) {
			out := []domaintoggle.ToggleWithDetail{{ID: "general/toggle1"}, {ID: "general/toggle2"}}
			return out, 2, nil
		}

		return s
	}

	t.Run("ok", func(t *testing.T) {
		s := setupTest()
		ctx := s.Context

		out, tot, err := s.ServiceStd.DoListToggles(ctx, servicetoggle.ParamsDoListToggles{})

		assert.NoError(t, err)
		assert.Equal(t, int64(2), tot)
		assert.NotEmpty(t, out)
	})

	t.Run("err on bad params", func(t *testing.T) {
		s := setupTest()

		ctx := s.Context
		_, _, err := s.ServiceStd.DoListToggles(ctx, servicetoggle.ParamsDoListToggles{
			Limit: 101,
		})

		assert.Error(t, err)
	})

	t.Run("err on fetch failed", func(t *testing.T) {
		s := setupTest()

		s.MockToggleStore.FetchPagedFunc = func(ctx context.Context, in storetoggle.ParamsFetchPaged) ([]domaintoggle.ToggleWithDetail, int64, error) {
			return nil, 0, assert.AnError
		}

		ctx := s.Context
		_, _, err := s.ServiceStd.DoListToggles(ctx, servicetoggle.ParamsDoListToggles{})

		assert.Error(t, err)
	})
}

func TestServiceStd_DoListStrayToggles(t *testing.T) {
	t.SkipNow()
}

func TestServiceStd_DoGetToggle(t *testing.T) {
	t.SkipNow()
}

func TestServiceStd_DoCreateToggle(t *testing.T) {
	t.SkipNow()
}

func TestServiceStd_DoRemoveToggle(t *testing.T) {
	toggleid := "general/toggle1"
	setupTest := func() *TestSuite {
		s := SetupSuite(t)

		s.MockToggleStore.FetchPagedFunc = func(ctx context.Context, in storetoggle.ParamsFetchPaged) ([]domaintoggle.ToggleWithDetail, int64, error) {
			out := []domaintoggle.ToggleWithDetail{{ID: "general/toggle1"}}
			return out, 2, nil
		}

		return s
	}

	t.Run("ok", func(t *testing.T) {
		s := setupTest()

		ctx := s.Context
		out, err := s.ServiceStd.DoRemoveToggle(ctx, toggleid)

		assert.NoError(t, err)
		assert.NotEmpty(t, out)

		assert.Equal(t, toggleid, out.ID)
	})

	t.Run("case error on check failure", func(t *testing.T) {
		s := setupTest()

		s.MockToggleStore.FetchPagedFunc = func(ctx context.Context, in storetoggle.ParamsFetchPaged) ([]domaintoggle.ToggleWithDetail, int64, error) {
			return nil, 0, assert.AnError
		}

		ctx := s.Context
		_, err := s.ServiceStd.DoRemoveToggle(ctx, toggleid)

		assert.Error(t, err)
	})

	t.Run("case err on fetch paged return nothing", func(t *testing.T) {
		s := setupTest()

		s.MockToggleStore.FetchPagedFunc = func(ctx context.Context, in storetoggle.ParamsFetchPaged) ([]domaintoggle.ToggleWithDetail, int64, error) {
			return nil, 0, nil
		}

		ctx := s.Context
		_, err := s.ServiceStd.DoRemoveToggle(ctx, toggleid)

		assert.Error(t, err)
	})

	t.Run("case err on remove toggles by ids return error", func(t *testing.T) {
		s := setupTest()

		s.MockToggleStore.RemoveTogglesByIDsFunc = func(ctx context.Context, ids []string) error {
			return assert.AnError
		}

		ctx := s.Context
		_, err := s.ServiceStd.DoRemoveToggle(ctx, toggleid)

		assert.Error(t, err)
	})

	t.Run("case err on delete access log by toggle id return error", func(t *testing.T) {
		s := setupTest()

		s.MockAccessLogService.DeleteAccessLogByToggleIDFunc = func(ctx context.Context, toggleid string) error {
			return assert.AnError
		}

		ctx := s.Context
		_, err := s.ServiceStd.DoRemoveToggle(ctx, toggleid)

		assert.Error(t, err)
	})

	t.Run("case err on saga rollback failed", func(t *testing.T) {
		s := setupTest()

		ctx := s.Context
		ctxsaga.CreateSaga(ctx)

		s.MockAccessLogService.DeleteAccessLogByToggleIDFunc = func(ctx context.Context, toggleid string) error {
			saga, ok := ctxsaga.GetSaga(ctx)
			require.True(t, ok)
			saga.AddRollbackFx(func() error { return assert.AnError })
			return assert.AnError
		}

		_, err := s.ServiceStd.DoRemoveToggle(ctx, toggleid)

		assert.Error(t, err)
	})

	t.Run("case err on saga commit failed", func(t *testing.T) {
		s := setupTest()

		ctx := s.Context
		ctxsaga.CreateSaga(ctx)

		s.MockAccessLogService.DeleteAccessLogByToggleIDFunc = func(ctx context.Context, toggleid string) error {
			saga, ok := ctxsaga.GetSaga(ctx)
			require.True(t, ok)
			saga.AddCommitFx(func() error { return assert.AnError })
			return nil
		}

		_, err := s.ServiceStd.DoRemoveToggle(ctx, toggleid)

		assert.Error(t, err)
	})

}

func TestServiceStd_DoStatToggle(t *testing.T) {
	t.SkipNow()
}
