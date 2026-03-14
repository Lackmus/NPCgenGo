# NPCgenGo â€” running and data directory

Quick notes on running the binary and where it looks for `data/`.

## Project structure

Core layout:

- `cmd/` â€” executable entrypoints and CLI wiring
- `internal/app/` — application orchestration layer (controllers, web adapters, views)
- `pkg/` â€” reusable core application logic (`pkg/...`)
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
- `cmd/npcgen-mapper/main.go` â€” JSON mapper demo/validation CLI

## Architecture walkthrough

Typical flow (CLI or HTTP):

1. Entry point: `cmd/*/main.go` binary starts
2. Root facade: `npcgen.go` delegates to `internal/app/npcgen_app.go`
3. App layer: wires `internal/app/controllers` (UI-agnostic) and `ui/console` views
4. Web layer: `internal/app/web` HTTP server routes to controllers
5. Controllers: call `pkg/service` domain services
6. Services: use `pkg/model` and `pkg/shared` contracts
7. Persistence: `internal/platform/loader` adapters handle config/storage
8. Data: read/write in `data/creation_data` and `data/npc_database`

Quick mental model:

- `cmd` = startup and process lifecycle
- `internal/app/controllers` = UI-agnostic orchestration (shared by all UIs)
- `internal/app/web` = HTTP-specific adapters (routes, middleware)
- `internal/app` = app wiring and initialization
- `pkg` = reusable business/domain logic
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

## Using as a dependency

`npcgengo.NewNPCGen()` reads JSON from runtime filesystem paths, not from embedded creation data.
When consuming this module from another repo/app, pass an explicit data directory.

```go
package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/lackmus/npcgengo"
)

func main() {
	// Optional override (works well in CI/container/dev)
	dataDir := os.Getenv("NPCGEN_DATA")
	if dataDir == "" {
		// Cross-platform default relative to your app's working directory.
		// Adjust this to where you copy/provision npcgen data in your app.
		dataDir = filepath.Join("assets", "npcgen", "data")
	}

	absDataDir, err := filepath.Abs(dataDir)
	if err != nil {
		log.Fatal(err)
	}

	app, err := npcgengo.NewNPCGenWithDataDir(absDataDir)
	if err != nil {
		log.Fatal(err)
	}

	_ = app
}
```

Tip: the provided directory should contain `creation_data/` and `npc_database/`.

## Data quality check

Before pushing JSON changes, run:

```bash
go run ./cmd/npcgen-doctor --data-dir ./data
```

What it validates:
- faction → species references
- species → name source references
- trait `Opposes` references
- NPC type/subtype consistency and required subtype fields (stats, equipment options)

CI runs the same command, so invalid creation data will fail pull requests.

## Run modes

| Mode | Entry point | Command | Notes |
|---|---|---|---|
| Web/API host | `cmd/npcgen-web/main.go` | `go run ./cmd/npcgen-web --data-dir "G:\My Drive\RootProject\NPCgenGo\data"` | Serves API and web UI at `/ui/` |
| Console UI | `cmd/npcgen-console/main.go` | `go run ./cmd/npcgen-console --data-dir "G:\My Drive\RootProject\NPCgenGo\data"` | Interactive terminal UI |
| Desktop UI (Wails) | `cmd/npcgen-wails/main.go` | `go run ./cmd/npcgen-wails --data-dir "G:\\My Drive\\RootProject\\NPCgenGo\\data"` | Wails desktop host |
| Mapper demo | `cmd/npcgen-mapper/main.go` | `go run ./cmd/npcgen-mapper --data-dir ./data --in ./npc.json` | Round-trips `pkg/mapper.NPCInput` JSON through the public mapper and validates it |

Mapper demo input example:

```json
{
	"id": "npc-42",
	"name": "Alice Smith",
	"type": "Civilian",
	"subtype": "someCivilianSubtypeID",
	"species": "someSpeciesID",
	"faction": "someFactionID",
	"trait": "someTraitID",
	"stats": "STR:2, DEX:1",
	"items": "Fists",
	"notes": "Scouting routes near river crossing"
}
```

You can also pipe JSON directly:

```powershell
Get-Content .\npc.json | go run ./cmd/npcgen-mapper --data-dir .\data
```

Running the example host (same wiring):
```powershell
cd examples/embed
go run . --data-dir "..\..\data"
```

Notes
- For distribution, you can place the `data/` directory next to the installed binary and run without flags.
- If you want a single-file release, consider adding a build mode that uses Go `embed` to include `data/` at build-time.
