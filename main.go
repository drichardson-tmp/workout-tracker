package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"workout-tracker/backend"
	"workout-tracker/backend/db"
	"workout-tracker/backend/middleware"
)

func main() {
	// Load .env if present; silently ignored when absent (e.g. production).
	_ = godotenv.Load()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	database, err := db.Connect(dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Zitadel auth is optional: set ZITADEL_DOMAIN to enable it.
	var authorizer *middleware.Authorizer
	if domain := os.Getenv("ZITADEL_DOMAIN"); domain != "" {
		authorizer, err = middleware.NewAuthorizer(
			context.Background(),
			domain,
			os.Getenv("ZITADEL_PORT"),
			os.Getenv("ZITADEL_CLIENT_ID"),
			os.Getenv("ZITADEL_CLIENT_SECRET"),
		)
		if err != nil {
			log.Fatal("Failed to create Zitadel authorizer:", err)
		}
	}

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	r.Use(middleware.Auth(authorizer))

	config := huma.DefaultConfig("Workout Tracker API", "1.0.0")
	api := humagin.New(r, config)

	backend.RegisterRoutes(api, database)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server listening on :%s", port)
	log.Printf("API docs: http://localhost:%s/docs", port)
	log.Printf("OpenAPI spec: http://localhost:%s/openapi.json", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
