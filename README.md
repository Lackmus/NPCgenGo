# NPCgenGo â€” running and data directory

Quick notes on running the binary and where it looks for `data/`.

## Project structure

Core layout:

- `cmd/` â€” executable entrypoints and CLI wiring
- `internal/app/` — application orchestration layer (controllers, web adapters, views)
- `pkg/` â€” reusable core application logic (`pkg/product/...`)
- `internal/` â€” module-private implementation details (`internal/platform/...`)

`cmd/` is entrypoint-only (`cmd/npcgen-web/main.go`, `cmd/npcgen-console/main.go`, `cmd/npcgen-wails/main.go`). App wiring lives in `internal/app/`.

UI and runtime data:
- `data/` â€” runtime JSON creation data and local NPC database files
- `ui/web/` â€” browser UI assets
- `ui/console/` â€” console UI assets (Go package)
- `ui/fyne/` â€” desktop UI notes
- `ui/wails/` â€” Wails frontend assets for desktop app

Entrypoints:
- `cmd/npcgen-web/main.go` â€” web/API host
- `cmd/npcgen-console/main.go` â€” console host
- `cmd/npcgen-wails/main.go` â€” desktop UI host (Wails)

## Architecture walkthrough

Typical flow (CLI or HTTP):

1. Entry point: `cmd/*/main.go` binary starts
2. Root facade: `NPCGen.go` delegates to `internal/app/npcgen_app.go`
3. App layer: wires `internal/app/controllers` (UI-agnostic) and `ui/console` views
4. Web layer: `internal/app/web` HTTP server routes to controllers
5. Controllers: call `pkg/product/service` domain services
6. Services: use `pkg/product/model` and `pkg/product/shared` contracts
7. Persistence: `internal/platform/loader` adapters handle config/storage
8. Data: read/write in `data/creation_data` and `data/npc_database`

Quick mental model:

- `cmd` = startup and process lifecycle
- `internal/app/controllers` = UI-agnostic orchestration (shared by all UIs)
- `internal/app/web` = HTTP-specific adapters (routes, middleware)
- `internal/app` = app wiring and initialization
- `pkg/product` = reusable business/domain logic
- `internal/platform` = infrastructure/plumbing (JSON loaders, helpers)
- `data` = runtime content

Priority order for data directory resolution (implemented):

- `--data-dir <path>` CLI flag (highest precedence)
- `NPCGEN_DATA` environment variable
- Auto-discovery by searching upward from the current working directory for `data/creation_data/factiondata`
- If still not found, auto-discovery by searching upward from the executable directory
- Final fallback: literal `data` relative path

Examples:
- Running `wails build` from `cmd/npcgen-wails` still resolves project-root `data/` automatically.
- Passing `--data-dir` to either a project root or an explicit `.../data` path is supported.

## Run modes

| Mode | Entry point | Command | Notes |
|---|---|---|---|
| Web/API host | `cmd/npcgen-web/main.go` | `go run ./cmd/npcgen-web --data-dir "G:\My Drive\RootProject\NPCgenGo\data"` | Serves API and web UI at `/ui/` |
| Console UI | `cmd/npcgen-console/main.go` | `go run ./cmd/npcgen-console --data-dir "G:\My Drive\RootProject\NPCgenGo\data"` | Interactive terminal UI |
| Desktop UI (Wails) | `cmd/npcgen-wails/main.go` | `go run ./cmd/npcgen-wails --data-dir "G:\\My Drive\\RootProject\\NPCgenGo\\data"` | Wails desktop host |

Running the example host (same wiring):
```powershell
cd examples/embed
go run . --data-dir "..\..\data"
```

Notes
- For distribution, you can place the `data/` directory next to the installed binary and run without flags.
- If you want a single-file release, consider adding a build mode that uses Go `embed` to include `data/` at build-time.
