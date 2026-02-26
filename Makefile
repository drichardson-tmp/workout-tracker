.PHONY: help install dev dev-backend dev-frontend generate migrate-new migrate-up migrate-status build test clean lint-frontend format-frontend db-start db-stop db-reset db-logs auth-start auth-stop auth-logs

help: ## Show this help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-25s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# ── dependencies ──────────────────────────────────────────────────────────────

install: ## Install all dependencies (Go + bun) and configure git hooks
	go mod download
	cd frontend && bun install
	git config core.hooksPath .githooks

install-atlas: ## Install the Atlas CLI (macOS / Linux)
	curl -sSf https://atlasgo.sh | sh

# ── database (Docker, persistent volume) ──────────────────────────────────────
# Uses a named Docker volume so data survives container restarts.
# workout_tracker_dev is a scratch DB used only by Atlas during migrate-new.

DB_CONTAINER  = workout-tracker-db
DB_VOLUME     = workout-tracker-pgdata
DB_USER       = postgres
DB_PASS       = postgres
DB_NAME       = workout_tracker
DB_DEV_NAME   = workout_tracker_dev
DB_PORT       = 5432
PG_IMAGE      = postgres:18

db-start: ## Start Postgres container (creates volume if needed)
	@if [ "$$(docker ps -q -f name=^$(DB_CONTAINER)$$)" ]; then \
	  echo "$(DB_CONTAINER) is already running"; \
	elif [ "$$(docker ps -aq -f name=^$(DB_CONTAINER)$$)" ]; then \
	  echo "Restarting stopped container $(DB_CONTAINER)..."; \
	  docker start $(DB_CONTAINER); \
	  until docker exec $(DB_CONTAINER) pg_isready -U $(DB_USER) -q; do sleep 1; done; \
	  docker exec $(DB_CONTAINER) createdb -U $(DB_USER) $(DB_DEV_NAME) 2>/dev/null || true; \
	  echo "Postgres ready — $(DB_NAME) and $(DB_DEV_NAME) available on :$(DB_PORT)"; \
	else \
	  docker run -d \
	    --name $(DB_CONTAINER) \
	    -e POSTGRES_USER=$(DB_USER) \
	    -e POSTGRES_PASSWORD=$(DB_PASS) \
	    -e POSTGRES_DB=$(DB_NAME) \
	    -p $(DB_PORT):5432 \
	    -v $(DB_VOLUME):/var/lib/postgresql \
	    $(PG_IMAGE); \
	  echo "Waiting for Postgres to be ready..."; \
	  until docker exec $(DB_CONTAINER) pg_isready -U $(DB_USER) -q; do sleep 1; done; \
	  docker exec $(DB_CONTAINER) createdb -U $(DB_USER) $(DB_DEV_NAME) 2>/dev/null || true; \
	  echo "Postgres ready — $(DB_NAME) and $(DB_DEV_NAME) available on :$(DB_PORT)"; \
	fi

db-stop: ## Stop (but keep) the Postgres container
	docker stop $(DB_CONTAINER) && docker rm $(DB_CONTAINER)

db-reset: ## Destroy all data and recreate a fresh database (WARNING: data loss)
	-docker stop $(DB_CONTAINER) 2>/dev/null; docker rm $(DB_CONTAINER) 2>/dev/null
	docker volume rm $(DB_VOLUME) 2>/dev/null || true
	$(MAKE) db-start

db-logs: ## Tail Postgres container logs
	docker logs -f $(DB_CONTAINER)

# ── Zitadel auth (docker compose) ─────────────────────────────────────────────
# Spins up Zitadel and its own Postgres on port 8081.
# After first start, open http://localhost:8081 and complete initial setup,
# then create an API application and set ZITADEL_* vars in your .env.

AUTH_COMPOSE = docker compose -f docker-compose.auth.yml

auth-start: ## Start Zitadel + its Postgres (initialises on first run)
	$(AUTH_COMPOSE) up -d

auth-stop: ## Stop and remove Zitadel containers
	$(AUTH_COMPOSE) down

auth-logs: ## Tail Zitadel logs
	$(AUTH_COMPOSE) logs -f zitadel

# ── local dev ─────────────────────────────────────────────────────────────────

dev: ## Run backend and frontend concurrently (Ctrl-C stops both)
	@trap 'kill 0' INT TERM EXIT; \
	go run . & \
	(cd frontend && bun run dev) & \
	wait

dev-backend: ## Run the Go server only
	go run .

dev-frontend: ## Run the Vite dev server only
	cd frontend && bun run dev

# ── code generation ───────────────────────────────────────────────────────────

generate: ## Generate openapi.json → frontend/src/generated/api.ts
	go run ./cmd/genschema > openapi.json
	cd frontend && bun run generate

# ── database migrations (Atlas) ───────────────────────────────────────────────
# Requires: atlas CLI installed and a local Postgres running.
# The "dev" DB in atlas.hcl is used only as a scratch space by Atlas — it must
# exist but Atlas will manage its schema automatically.

migrate-new: ## Diff GORM models → new migration file  (usage: make migrate-new name=add_foo)
	atlas migrate diff $(name) --env local

migrate-up: ## Apply all pending migrations to the local DB
	atlas migrate apply --env local --url "postgres://$(DB_USER):$(DB_PASS)@localhost:$(DB_PORT)/$(DB_NAME)?sslmode=disable"

migrate-status: ## Show applied/pending migration status
	atlas migrate status --env local

# ── build & test ──────────────────────────────────────────────────────────────

build: ## Build the backend binary
	go build -o bin/server .

test: ## Run Go tests
	go test -v ./...

clean: ## Remove build artifacts
	rm -rf bin/ openapi.json
	cd frontend && rm -rf dist

# ── frontend quality ──────────────────────────────────────────────────────────

lint-frontend: ## Lint frontend with Biome
	cd frontend && bun run lint

format-frontend: ## Format frontend with Biome
	cd frontend && bun run format
