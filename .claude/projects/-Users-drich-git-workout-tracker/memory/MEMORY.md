# Project Memory

## Stack
- Go 1.25, module name: `workout-tracker`
- Router: gin-gonic/gin
- ORM: gorm.io/gorm + gorm.io/driver/postgres (pgx under the hood, no lib/pq)
- API/OpenAPI: github.com/danielgtaylor/huma/v2 v2.32.0
  - Gin adapter import: `github.com/danielgtaylor/huma/v2/adapters/humagin`
  - (Adapters moved to `adapters/` subdir as of ~v2.30; the old `humagin` top-level path no longer exists)
- DB migrations: Atlas CLI with ariga.io/atlas-provider-gorm v0.6.0
- Frontend: Vue 3 + Vite (port 3000) + Biome linter + openapi-typescript

## Directory layout
```
main.go                  # entry point
backend/
  db/db.go               # gorm.Open helper
  models/user.go         # GORM model
  models/workout.go      # GORM model
  schemas/user.go        # Huma request/response types
  schemas/workout.go
  handlers/user.go       # handler structs (NewXHandler(db))
  handlers/workout.go
  router.go              # RegisterRoutes(api, db) — package backend
cmd/
  atlasloader/main.go    # atlas.hcl `external_schema` loader
  genschema/main.go      # go run ./cmd/genschema > openapi.json
atlas.hcl                # Atlas config, env "local"
migrations/              # Atlas-managed SQL migration files
frontend/src/generated/  # openapi-typescript output (gitignored)
```

## Testing
- Test files: `backend/*_test.go`, package `backend_test` (external — avoids import cycle)
- Helper: `newMockDB(t)` → GORM + go-sqlmock (WithoutReturning:true for simpler Exec mocks)
- Helper: `newTestAPI(t, db)` → humatest.TestAPI with all routes registered
- Deps: `github.com/DATA-DOG/go-sqlmock v1.5.2`, `github.com/stretchr/testify`
- Run: `make test` or `go test ./backend/...`

## Key commands
- `make dev` — run backend + frontend concurrently
- `make generate` — gen openapi.json then frontend types
- `make migrate-new name=foo` — diff GORM models → new migration
- `make migrate-up` — apply migrations to local DB
