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

func TestListUsers_Empty(t *testing.T) {
	db, mock := newMockDB(t)
	api := newTestAPI(t, db)

	mock.ExpectQuery(`SELECT \* FROM "users"`).
		WillReturnRows(sqlmock.NewRows(userCols()))

	resp := api.Get("/api/users")

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.JSONEq(t, `[]`, resp.Body.String())
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestListUsers_WithResults(t *testing.T) {
	db, mock := newMockDB(t)
	api := newTestAPI(t, db)

	rows := sqlmock.NewRows(userCols()).
		AddRow(int64(1), fixedTime, fixedTime, nil, nil, "alice@example.com", "Alice", "hash1").
		AddRow(int64(2), fixedTime, fixedTime, nil, nil, "bob@example.com", "Bob", "hash2")
	mock.ExpectQuery(`SELECT \* FROM "users"`).WillReturnRows(rows)

	resp := api.Get("/api/users")

	require.Equal(t, http.StatusOK, resp.Code)
	var body []schemas.UserResponse
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &body))
	assert.Len(t, body, 2)
	assert.Equal(t, "alice@example.com", body[0].Email)
	assert.Equal(t, "Bob", body[1].Name)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUser_Found(t *testing.T) {
	db, mock := newMockDB(t)
	api := newTestAPI(t, db)

	mock.ExpectQuery(`SELECT \* FROM "users"`).
		WillReturnRows(sqlmock.NewRows(userCols()).
			AddRow(int64(1), fixedTime, fixedTime, nil, nil, "alice@example.com", "Alice", "hash"))

	resp := api.Get("/api/users/1")

	require.Equal(t, http.StatusOK, resp.Code)
	var body schemas.UserResponse
	require.NoError(t, json.Unmarshal(resp.Body.Bytes(), &body))
	assert.Equal(t, int64(1), body.ID)
	assert.Equal(t, "alice@example.com", body.Email)
	assert.Equal(t, "Alice", body.Name)
	// PasswordHash must NOT be exposed
	assert.NotContains(t, resp.Body.String(), "password")
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUser_NotFound(t *testing.T) {
	db, mock := newMockDB(t)
	api := newTestAPI(t, db)

	mock.ExpectQuery(`SELECT \* FROM "users"`).
		WillReturnRows(sqlmock.NewRows(userCols())) // no rows → 404

	resp := api.Get("/api/users/99")

	assert.Equal(t, http.StatusNotFound, resp.Code)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateUser(t *testing.T) {
	db, mock := newMockDB(t)
	api := newTestAPI(t, db)

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO "users"`).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	resp := api.Post("/api/users", map[string]any{
		"email":    "alice@example.com",
		"name":     "Alice",
		"password": "supersecret",
	})

	assert.Equal(t, http.StatusCreated, resp.Code)
	// Password must not leak into the response body
	assert.NotContains(t, resp.Body.String(), "supersecret")
	require.NoError(t, mock.ExpectationsWereMet())
}
