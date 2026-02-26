# Workout Tracker

A small personal project for learning and experimenting with a specific modern Go + Vue stack. Tooling choices prioritise fast iteration over production hardening (or truly safe user monitoring, having chosen self-hosted Zitadel).

## Tech Stack

### Backend
- **Go 1.25** — application server
- **Gin** — HTTP router
- **Huma v2** — OpenAPI schema generation and request/response validation on top of Gin
- **GORM** — ORM for model definitions and query building
- **PostgreSQL** — primary database (run locally via Docker)
- **Atlas** — database migration diffing and application, driven from GORM models
- **Zitadel** — OIDC authentication and authorisation (token introspection)
- **godotenv** — `.env` file loading for local development

### Frontend
- **Vue 3** (`<script setup>`) — UI framework
- **Vite** — dev server and bundler
- **Vue Router** — client-side routing
- **Axios** — HTTP client
- **openapi-typescript** — generates a typed API client from the backend's `openapi.json`

### Tooling
- **Bun** — frontend package manager and script runner
- **Biome** — linter and formatter for the frontend (replaces ESLint + Prettier)
- **Docker** — runs Postgres and Zitadel locally
- **Make** — task runner (`make help` lists all targets)
- **Git hooks** (`.githooks/`) — on commit: `go build`, `go vet`, Biome, generated API freshness check, and Atlas migration drift check

## Getting Started

```sh
make install       # install Go + frontend deps, configure git hooks
make db-start      # start local Postgres in Docker
make migrate-up    # apply migrations
make dev           # run backend + frontend concurrently
```

Copy `.env.example` to `.env` and fill in values before running.
