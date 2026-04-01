# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Documentation-First Protocol

**Before starting any formal task**, read these files in order:
1. `agent.md` — master guide for this repo (primary authority)
2. `HANDOFF.md` — current batch task details
3. `docs/project/agent-execution-rules.md`
4. `docs/project/current-status-and-direction.md`
5. `docs/project/project-memory.md`
6. For frontend tasks, also read: `docs/frontend/overview.md`, `docs/frontend/issues.md`, `docs/frontend/specs/`, `docs/frontend/plans/`

Important conclusions must be written back to docs — not left only in chat.

## Project Structure

Source root is `unifocus/` inside the repo root `Unifocus-v1.2/`.

- `unifocus/backend/` — Go + Gin API (primary backend, port 8080)
- `unifocus/web/` — Next.js 14 frontend (primary production frontend, port 3000)
- `unifocus/web-vite/` — Vite + React experimental frontend (port 5173)
- `unifocus/nlp-service/` — Python FastAPI NLP microservice (port 8000)
- `unifocus/migrations/` — PostgreSQL migration SQL files
- `unifocus/Makefile` — all common dev commands

## Development Commands

All `make` commands run from `unifocus/`:

```bash
# Docker services (postgres, redis, api, nlp)
make up              # Start all services
make down            # Stop all services
make logs-api        # Tail API logs
make restart-api     # Restart Go API container

# Database
make db-migrate      # Run migration files
make db-reset        # Reset schema
make db-shell        # Interactive psql

# Local dev (outside Docker)
make dev-backend     # go run cmd/api/main.go
make dev-nlp         # uvicorn app.main:app --reload

# Frontend (run from unifocus/web/)
npm run dev          # Next.js dev server
npm run build        # Production build
npm run lint         # ESLint

# Code quality
make fmt-go          # gofmt
make lint-go         # go vet
make fmt-python      # black

# Health checks
make test-api        # curl localhost:8080/health
make test-nlp        # curl localhost:8000/health
```

## Architecture

### Frontend (unifocus/web/)

Next.js App Router. Key paths:
- `app/page.tsx` — redirects `/` → `/dashboard`
- `app/dashboard/page.tsx` — main dashboard
- `app/dashboard/planet/page.tsx` — Planet 3D view
- `components/dashboard/` — page-level business components
  - `planet-client.tsx` — Planet view orchestrator
  - `planet-scene.tsx` — Three.js 3D rendering
  - `planet-data.ts` — data transformation (proximity/overload detection)
- `components/ui/` — shadcn/ui primitives
- `lib/api/` — axios HTTP client

API calls rewrite `/api/*` → `$API_BASE_URL/api/*` (default: `http://localhost:8080`). The 401 handler does **not** redirect to `/login` (login not yet implemented).

Build artifacts: dev uses `.next-dev-{port}/`, prod uses `.next/`. If artifacts seem corrupted, delete the relevant directory and rebuild.

### Backend (unifocus/backend/)

Go + Gin. Entry: `cmd/api/main.go`.

Layer structure:
- `internal/api/handlers/` — HTTP handlers
- `internal/api/middleware/` — CORS, auth JWT
- `internal/service/` — business logic
- `internal/repository/postgres/` — SQL queries
- `internal/repository/redis/` — cache
- `internal/crawler/scrapers/` — web scraping

### NLP Service (unifocus/nlp-service/)

FastAPI + Python. Not a blocking dependency in current dev phase.

### Data

- PostgreSQL (port 5432) + Redis (port 6379) via Docker
- Schema managed by `migrations/` SQL files
- Opportunity data may return `null` from API — always guard with empty array fallback at both data and render layers

## Frontend Debugging Checklist

When a page is inaccessible or crashes, check in this order:
1. Is the correct `next dev` process running on the expected port? (Kill stale processes)
2. Are `.next` artifacts corrupted or mixed with a build output? (Delete and restart)
3. Is the route actually defined in `app/`?
4. Is the API returning `null` where an array is expected?
5. Are there sidebar links pointing to unimplemented routes?

## Key Constraints

- `unifocus/web-vite/` is an experimental parallel line — do not merge concerns with `unifocus/web/`
- `nlp-service` is non-blocking for current phase
- Frontend tasks must confirm real routes, real API contracts, and real auth flow before changing code
