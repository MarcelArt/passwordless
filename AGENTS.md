# AGENTS.md

Guidance for agents working in this repo.

## Session rules (must follow)

- **Do not run the server or migrations.** The user runs these themselves. Avoid `make dev`, `make go`, `make migrate`, `make migrate-force`, `docker-run`, `compose`, and any `go run main.go serve|migrate`. Safe verification only: `go build ./...`, `go vet ./...`, and `make swag` (docs generation) if needed.
- **Stay on your side of the repo.** `internal/`, `cmd/`, `pkg/`, `main.go` are the server side. `web/` is the UI side. If the current task is on the server side, do **not** edit `web/` (and vice versa). If a change genuinely requires touching the other side, stop and ask the user instead of editing it.

## Architecture

Single Go binary (Gin) that serves **both** the JSON API and the server-rendered UI — `web/` is not a separate frontend app. It is Go `html/template` + HTMX rendered by the same process and consumes `internal/` code.

- Entry: `main.go` → `cmd.Execute()` (Cobra CLI). Subcommands: `serve`, `migrate [--drop]`.
- API wiring: `cmd/serve.go` → `internal/v1/routes/route.go` (group `/api/v1`) + `web/routes/web.go`.
- Layered backend in `internal/v1/`: `routes` → `handlers` → `services` → `repositories` → `entities`. Follow this direction; don't skip layers.
- Shared helpers: `internal/common/` (CRUD base, JWT, cookie, WebAuthn, response, permission). Config: `internal/configs/` (`env.go`, `postgres.go`, `oauth2.go`). GORM models: `internal/entities/`.
- Handlers are registered with constructor injection (repo → service → handler); mirror this when adding endpoints.
- Auth is passwordless: WebAuthn + JWT, plus OAuth2. WebAuthn config (RP_ID, RP_ORIGINS, RP_DISPLAY_NAME) is read in `internal/configs/env.go`.

## Critical build quirk: `docs/` is generated and gitignored

`docs/*` is in `.gitignore`, but `main.go` does a blank import of `github.com/MarcelArt/passwordless/docs`. If `docs/` is missing, **the project will not compile**.

- Regenerate with `make swag` (i.e. `swag init --parseDependency --parseInternal`) whenever Swagger annotations in handlers/routes change or before building from a fresh checkout.
- `docs/docs.go` is auto-generated ("DO NOT EDIT") — never hand-edit it. Edit the `@*` swag annotations in source instead.

## Database

- No SQL migration files. Schema is managed via GORM `AutoMigrate` in `internal/configs/postgres.go`, run through the `migrate` command (which the user runs, not the agent).
- To add/alter a model: update `internal/entities/` and ensure it is listed in `MigrateDB()`.

## Configuration

- Config comes from a gitignored `.env` loaded by godotenv (`internal/configs/env.go`). Copy `example.env` to `.env`.
- Note: `example.env` is **incomplete** — it omits the WebAuthn vars (`RP_ID`, `RP_ORIGINS`, `RP_DISPLAY_NAME`) that `env.go` reads. Flag this if the user hits empty-config issues.

## Conventions / gotchas

- No tests, lint config, CI workflows, or README exist in this repo. Do not invent test/lint commands; ask the user if verification is needed.
- Root Cobra command `Use:` is `"lepas"` and the Docker binary is named `lepas` — the project is sometimes referred to as "lepas" despite the repo name "passwordless".
- `tmp/` (air live-reload temp) is gitignored.
- When editing UI (`web/`), templates are loaded by filesystem path at runtime (`web/views/layout.html` + page template in `web/handlers/render.go`); the Docker image copies `web/` into the binary's working dir.
