package backend

import (
	"context"
	"net/http"

	"workout-tracker/backend/handlers"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
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
	uh.RegisterRoutes(api)
	wh.RegisterRoutes(api)
}
