// genschema generates openapi.json from Huma route registrations.
// Run via: go run ./cmd/genschema > openapi.json
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-gonic/gin"
	"workout-tracker/backend"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	config := huma.DefaultConfig("Workout Tracker API", "1.0.0")
	api := humagin.New(r, config)

	// Register routes with a nil DB — Huma introspects types at registration
	// time and never invokes handlers during schema generation.
	backend.RegisterRoutes(api, nil)

	b, err := json.MarshalIndent(api.OpenAPI(), "", "  ")
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to marshal OpenAPI spec:", err)
		os.Exit(1)
	}
	fmt.Println(string(b))
}
