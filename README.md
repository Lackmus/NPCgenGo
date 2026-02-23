# NPCgenGo â€” running and data directory

Quick notes on running the binary and where it looks for `data/`.

## Project structure

Applied structure scope (only):

- `cmd/` â€” executable entrypoints and CLI wiring
- `internal/app/` â€” application orchestration layer (controllers, server adapters, views)
- `pkg/` â€” reusable core application logic (`pkg/product/...`)
- `internal/` â€” module-private implementation details (`internal/platform/...`)

No additional optional folders from that setup were added.

`cmd/` is now kept entrypoint-only (`cmd/npcgen/main.go`), while app wiring lives in `internal/app/`.

Current project folders also include:
- `data/` â€” runtime JSON creation data and local NPC database files
- `web_demo/` â€” web demonstration assets

## Architecture walkthrough

Typical flow (CLI or HTTP):

1. Entry point starts in `cmd/npcgen/main.go`.
2. Root facade (`NPCGen.go`) delegates app setup to `internal/app/npcgen_app.go`.
3. App layer wires handlers/controllers/views in `internal/app/handlers` and `internal/app/view`.
4. Handlers call domain services in `pkg/product/service`.
5. Services use domain models/contracts in `pkg/product/model` and `pkg/product/shared`.
6. Persistence/config loading is handled by platform adapters in `internal/platform/loader`.
7. Data is read/written in `data/creation_data` and `data/npc_database`.

Quick mental model:

- `cmd` = startup and process lifecycle
- `internal/app` = orchestration and adapters (HTTP/controller/view wiring)
- `pkg/product` = reusable business/domain logic
- `internal/platform` = infrastructure/plumbing (JSON loaders, helpers)
- `data` = runtime content

Priority order for data directory resolution (implemented):

- `--data-dir <path>` CLI flag (highest precedence)
- `NPCGEN_DATA` environment variable
- `data/` directory located next to the running executable (useful for installed binaries)
- `./data` in the current working directory (developer-friendly default)

Developer examples

PowerShell (flag):
```powershell
go run ./cmd/npcgen --data-dir "G:\My Drive\RootProject\NPCgenGo\data"
```

PowerShell (env var):
```powershell
$env:NPCGEN_DATA = 'G:\My Drive\RootProject\NPCgenGo\data'
go run ./cmd/npcgen
```

Running the example host (same wiring):
```powershell
cd examples/embed
go run . --data-dir "..\..\data"
```

Notes
- For distribution, you can place the `data/` directory next to the installed binary and run without flags.
- If you want a single-file release, consider adding a build mode that uses Go `embed` to include `data/` at build-time.
