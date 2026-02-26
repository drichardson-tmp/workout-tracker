package backend_test

import (
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"workout-tracker/backend"
)

// fixedTime is a stable timestamp used across test assertions.
var fixedTime = time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

// newMockDB returns a GORM DB backed by go-sqlmock and its companion mock
// controller. WithoutReturning is set so that INSERT/UPDATE/DELETE produce
// Exec calls (easier to mock) rather than Query…RETURNING calls.
func newMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()
	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn:             sqlDB,
		WithoutReturning: true,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)

	t.Cleanup(func() { sqlDB.Close() })
	return db, mock
}

// newTestAPI wires all routes onto a humatest API backed by db.
func newTestAPI(t *testing.T, db *gorm.DB) humatest.TestAPI {
	t.Helper()
	_, api := humatest.New(t)
	backend.RegisterRoutes(api, db)
	return api
}

// workoutCols returns the column names that GORM scans for a Workout row.
func workoutCols() []string {
	return []string{"id", "created_at", "updated_at", "deleted_at",
		"user_id", "name", "description", "duration_minutes"}
}

// userCols returns the column names that GORM scans for a User row.
func userCols() []string {
	return []string{"id", "created_at", "updated_at", "deleted_at",
		"zitadel_id", "email", "name", "password_hash"}
}
