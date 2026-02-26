package backend_test

import (
	"encoding/json"
	"net/http"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"workout-tracker/backend/schemas"
)

func TestListWorkouts_Empty(t *testing.T) {
	db, mock := newMockDB(t)
	api := newTestAPI(t, db)

	mock.ExpectQuery(`SELECT \* FROM "workouts"`).
		WillReturnRows(sqlmock.NewRows(workoutCols()))

	resp := api.Get("/api/workouts")

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.JSONEq(t, `[]`, resp.Body.String())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestListWorkouts_WithResults(t *testing.T) {
	db, mock := newMockDB(t)
	api := newTestAPI(t, db)

	rows := sqlmock.NewRows(workoutCols()).
		AddRow(int64(1), fixedTime, fixedTime, nil, int64(1), "Morning Run", "5km run", 30).
		AddRow(int64(2), fixedTime, fixedTime, nil, int64(1), "Evening Yoga", "", 45)
	mock.ExpectQuery(`SELECT \* FROM "workouts"`).WillReturnRows(rows)

	resp := api.Get("/api/workouts")

	require.Equal(t, http.StatusOK, resp.Code)
	var body []schemas.WorkoutResponse
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &body))
	assert.Len(t, body, 2)
	assert.Equal(t, "Morning Run", body[0].Name)
	assert.Equal(t, "Evening Yoga", body[1].Name)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestListWorkouts_FilterByUser(t *testing.T) {
	db, mock := newMockDB(t)
	api := newTestAPI(t, db)

	rows := sqlmock.NewRows(workoutCols()).
		AddRow(int64(1), fixedTime, fixedTime, nil, int64(42), "Pull Day", "", 60)
	mock.ExpectQuery(`SELECT \* FROM "workouts"`).
		WillReturnRows(rows)

	resp := api.Get("/api/workouts?userId=42")

	require.Equal(t, http.StatusOK, resp.Code)
	var body []schemas.WorkoutResponse
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &body))
	assert.Len(t, body, 1)
	assert.Equal(t, int64(42), body[0].UserID)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetWorkout_Found(t *testing.T) {
	db, mock := newMockDB(t)
	api := newTestAPI(t, db)

	mock.ExpectQuery(`SELECT \* FROM "workouts"`).
		WillReturnRows(sqlmock.NewRows(workoutCols()).
			AddRow(int64(1), fixedTime, fixedTime, nil, int64(1), "Morning Run", "5km run", 30))

	resp := api.Get("/api/workouts/1")

	require.Equal(t, http.StatusOK, resp.Code)
	var body schemas.WorkoutResponse
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &body))
	assert.Equal(t, int64(1), body.ID)
	assert.Equal(t, "Morning Run", body.Name)
	assert.Equal(t, 30, body.DurationMinutes)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetWorkout_NotFound(t *testing.T) {
	db, mock := newMockDB(t)
	api := newTestAPI(t, db)

	mock.ExpectQuery(`SELECT \* FROM "workouts"`).
		WillReturnRows(sqlmock.NewRows(workoutCols())) // no rows → ErrRecordNotFound

	resp := api.Get("/api/workouts/99")

	assert.Equal(t, http.StatusNotFound, resp.Code)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateWorkout(t *testing.T) {
	db, mock := newMockDB(t)
	api := newTestAPI(t, db)

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO "workouts"`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	resp := api.Post("/api/workouts", map[string]any{
		"user_id":          1,
		"name":             "Morning Run",
		"description":      "5km run",
		"duration_minutes": 30,
	})

	assert.Equal(t, http.StatusCreated, resp.Code)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateWorkout(t *testing.T) {
	db, mock := newMockDB(t)
	api := newTestAPI(t, db)

	mock.ExpectQuery(`SELECT \* FROM "workouts"`).
		WillReturnRows(sqlmock.NewRows(workoutCols()).
			AddRow(int64(1), fixedTime, fixedTime, nil, int64(1), "Old Name", "", 20))

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "workouts" SET`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	resp := api.Patch("/api/workouts/1", map[string]any{
		"name":             "New Name",
		"duration_minutes": 45,
	})

	assert.Equal(t, http.StatusOK, resp.Code)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteWorkout(t *testing.T) {
	db, mock := newMockDB(t)
	api := newTestAPI(t, db)

	// Soft-delete: GORM issues UPDATE SET deleted_at=... rather than DELETE.
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "workouts" SET "deleted_at"`).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	resp := api.Delete("/api/workouts/1")

	assert.Equal(t, http.StatusNoContent, resp.Code)
	require.NoError(t, mock.ExpectationsWereMet())
}
