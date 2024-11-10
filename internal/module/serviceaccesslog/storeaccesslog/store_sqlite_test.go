package storeaccesslog_test

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/avrebarra/goggle/internal/module/serviceaccesslog/domainaccesslog"
	"github.com/avrebarra/goggle/internal/module/serviceaccesslog/storeaccesslog"
	"github.com/avrebarra/goggle/utils/ctxboard"
	"github.com/avrebarra/goggle/utils/ctxsaga"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func normalize(in string) string { return regexp.QuoteMeta(in) }

type TestSuite struct {
	Context context.Context
	MockDB  sqlmock.Sqlmock
	Store   *storeaccesslog.StorageSQLite
}

func SetupSuite(t *testing.T) *TestSuite {
	s := &TestSuite{}

	gofakeit.Seed(333555444) // for deterministic tests

	s.Context = ctxboard.CreateWith(context.Background())

	db, mock, err := sqlmock.New()
	s.MockDB = mock
	require.NoError(t, err)

	s.MockDB.ExpectQuery("select sqlite_version()").WillReturnRows(sqlmock.NewRows([]string{"sqlite_version()"}).AddRow("3.36.1"))

	dialector := sqlite.New(sqlite.Config{Conn: db})
	gormdb, err := gorm.Open(dialector, &gorm.Config{})
	require.NoError(t, err)

	store, err := storeaccesslog.NewStorageSQLite(storeaccesslog.ConfigStorageSQLite{
		DB: gormdb,
	})
	require.NoError(t, err)
	s.Store = store

	return s
}

func TestNewStorageSQLite(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		s := SetupSuite(t)
		assert.NotNil(t, s.Store)
	})

	t.Run("on bad deps", func(t *testing.T) {
		_, err := storeaccesslog.NewStorageSQLite(storeaccesslog.ConfigStorageSQLite{
			DB: nil,
		})
		assert.Error(t, err)
	})
}

func TestStorageSQLite_FetchPaged(t *testing.T) {
	Q1 := normalize("SELECT count(*) FROM access_logs as log")
	Q2 := normalize("SELECT id, toggle_id, created_at FROM access_logs as log ORDER BY id asc LIMIT 10")
	Q3 := normalize("SELECT count(*) FROM access_logs as log WHERE log.toggle_id IN (?,?)")
	Q4 := normalize("SELECT id, toggle_id, created_at FROM access_logs as log WHERE log.toggle_id IN (?,?) ORDER BY id asc LIMIT 10")

	t.Run("ok", func(t *testing.T) {
		s := SetupSuite(t)

		s.MockDB.ExpectQuery(Q1).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))
		s.MockDB.ExpectQuery(Q2).WillReturnRows(sqlmock.NewRows([]string{"id", "toggle_id", "created_at"}).AddRow(1, "toggle1", "2024-11-09 22:27:54.798977+07:00").AddRow(2, "toggle2", "2024-11-09 22:27:54.798977+07:00"))

		ctx := s.Context
		out, tot, err := s.Store.FetchPaged(ctx, storeaccesslog.ParamsFetchPaged{})

		assert.NoError(t, err)
		assert.NotEmpty(t, out)
		assert.Equal(t, int64(2), tot)
	})

	t.Run("on invalid params", func(t *testing.T) {
		s := SetupSuite(t)

		ctx := s.Context
		_, _, err := s.Store.FetchPaged(ctx, storeaccesslog.ParamsFetchPaged{
			Offset: -5,
			Limit:  2000,
		})

		assert.Error(t, err)
	})

	t.Run("on filter by toggle id", func(t *testing.T) {
		s := SetupSuite(t)

		s.MockDB.ExpectQuery(Q3).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))
		s.MockDB.ExpectQuery(Q4).WillReturnRows(sqlmock.NewRows([]string{"id", "toggle_id", "created_at"}).AddRow(1, "toggle1", "2024-11-09 22:27:54.798977+07:00").AddRow(2, "toggle2", "2024-11-09 22:27:54.798977+07:00"))

		ctx := s.Context
		out, tot, err := s.Store.FetchPaged(ctx, storeaccesslog.ParamsFetchPaged{
			FilterToggleIDs: []string{"toggle1", "toggle2"},
		})

		assert.NoError(t, err)
		assert.NotEmpty(t, out)
		assert.Equal(t, int64(2), tot)
	})

	t.Run("on skip totaling", func(t *testing.T) {
		s := SetupSuite(t)

		s.MockDB.ExpectQuery(Q2).WillReturnRows(sqlmock.NewRows([]string{"id", "toggle_id", "created_at"}).AddRow(1, "toggle1", "2024-11-09 22:27:54.798977+07:00").AddRow(2, "toggle2", "2024-11-09 22:27:54.798977+07:00"))

		ctx := s.Context
		out, tot, err := s.Store.FetchPaged(ctx, storeaccesslog.ParamsFetchPaged{
			SkipTotal: true,
		})

		assert.NoError(t, err)
		assert.NotEmpty(t, out)
		assert.Equal(t, int64(0), tot)
	})

	t.Run("on fetch failed", func(t *testing.T) {
		s := SetupSuite(t)

		s.MockDB.ExpectQuery(Q1).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))
		s.MockDB.ExpectQuery(Q2).WillReturnError(assert.AnError)

		ctx := s.Context
		_, _, err := s.Store.FetchPaged(ctx, storeaccesslog.ParamsFetchPaged{})

		assert.Error(t, err)
	})
}

func TestStorageSQLite_CreateLog(t *testing.T) {
	Q1 := normalize("INSERT INTO `access_logs` (`toggle_id`,`created_at`,`id`) VALUES (?,?,?) RETURNING `id`")

	t.Run("ok", func(t *testing.T) {
		s := SetupSuite(t)

		accesslog := domainaccesslog.AccessLog{}.Fake()

		s.MockDB.ExpectBegin()
		s.MockDB.ExpectQuery(Q1).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(accesslog.ID))
		s.MockDB.ExpectCommit()

		ctx := s.Context
		saga := ctxsaga.CreateSaga(ctx)
		out, err := s.Store.CreateLog(ctx, accesslog)
		saga.Commit()

		assert.NoError(t, err)
		assert.NotEmpty(t, out)
	})

	t.Run("on bad data", func(t *testing.T) {
		s := SetupSuite(t)

		accesslog := domainaccesslog.AccessLog{}.Fake()
		accesslog.ToggleID = ""

		ctx := s.Context
		out, err := s.Store.CreateLog(ctx, accesslog)

		assert.Error(t, err)
		assert.Empty(t, out)
	})

	t.Run("on saga rollback", func(t *testing.T) {
		s := SetupSuite(t)

		accesslog := domainaccesslog.AccessLog{}.Fake()

		s.MockDB.ExpectBegin()
		s.MockDB.ExpectQuery(Q1).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(accesslog.ID))
		s.MockDB.ExpectRollback()

		ctx := s.Context
		saga := ctxsaga.CreateSaga(ctx)
		out, err := s.Store.CreateLog(ctx, accesslog)
		require.NoError(t, err)

		err = saga.Rollback()
		assert.NoError(t, err)
		assert.NotEmpty(t, out)
	})

	t.Run("on query failed", func(t *testing.T) {
		s := SetupSuite(t)

		accesslog := domainaccesslog.AccessLog{}.Fake()

		s.MockDB.ExpectBegin()
		s.MockDB.ExpectQuery(Q1).WillReturnError(assert.AnError)
		s.MockDB.ExpectRollback()

		ctx := s.Context
		saga := ctxsaga.CreateSaga(ctx)
		out, err := s.Store.CreateLog(ctx, accesslog)
		require.Error(t, err)

		err = saga.Rollback()
		assert.NoError(t, err)
		assert.Empty(t, out)
	})
}

func TestStorageSQLite_DeleteAllByToggleIDs(t *testing.T) { t.SkipNow() }
