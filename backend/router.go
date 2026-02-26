package backend

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
	"workout-tracker/backend/handlers"
)

type healthBody struct {
	Status string `json:"status"`
}

type healthOutput struct {
	Body healthBody
}

// RegisterRoutes wires all API routes onto the given Huma API.
// db may be nil when called from the schema generator (routes are registered
// for type introspection only; handlers are never invoked).
func RegisterRoutes(api huma.API, db *gorm.DB) {
	huma.Register(api, huma.Operation{
		OperationID: "health",
		Method:      http.MethodGet,
		Path:        "/health",
		Summary:     "Health check",
	}, func(_ context.Context, _ *struct{}) (*healthOutput, error) {
		return &healthOutput{Body: healthBody{Status: "ok"}}, nil
	})
	uh := handlers.NewUserHandler(db)
	wh := handlers.NewWorkoutHandler(db)

	// --- users ---
	huma.Register(api, huma.Operation{
		OperationID: "list-users",
		Method:      http.MethodGet,
		Path:        "/api/users",
		Summary:     "List all users",
		Tags:        []string{"users"},
	}, uh.ListUsers)

	huma.Register(api, huma.Operation{
		OperationID: "get-user",
		Method:      http.MethodGet,
		Path:        "/api/users/{userId}",
		Summary:     "Get a user by ID",
		Tags:        []string{"users"},
	}, uh.GetUser)

	huma.Register(api, huma.Operation{
		OperationID:   "create-user",
		Method:        http.MethodPost,
		Path:          "/api/users",
		Summary:       "Create a new user",
		Tags:          []string{"users"},
		DefaultStatus: http.StatusCreated,
	}, uh.CreateUser)

	// --- workouts ---
	huma.Register(api, huma.Operation{
		OperationID: "list-workouts",
		Method:      http.MethodGet,
		Path:        "/api/workouts",
		Summary:     "List workouts",
		Tags:        []string{"workouts"},
	}, wh.ListWorkouts)

	huma.Register(api, huma.Operation{
		OperationID: "get-workout",
		Method:      http.MethodGet,
		Path:        "/api/workouts/{workoutId}",
		Summary:     "Get a workout by ID",
		Tags:        []string{"workouts"},
	}, wh.GetWorkout)

	huma.Register(api, huma.Operation{
		OperationID:   "create-workout",
		Method:        http.MethodPost,
		Path:          "/api/workouts",
		Summary:       "Create a new workout",
		Tags:          []string{"workouts"},
		DefaultStatus: http.StatusCreated,
	}, wh.CreateWorkout)

	huma.Register(api, huma.Operation{
		OperationID: "update-workout",
		Method:      http.MethodPatch,
		Path:        "/api/workouts/{workoutId}",
		Summary:     "Update a workout",
		Tags:        []string{"workouts"},
	}, wh.UpdateWorkout)

	huma.Register(api, huma.Operation{
		OperationID:   "delete-workout",
		Method:        http.MethodDelete,
		Path:          "/api/workouts/{workoutId}",
		Summary:       "Delete a workout",
		Tags:          []string{"workouts"},
		DefaultStatus: http.StatusNoContent,
	}, wh.DeleteWorkout)
}
