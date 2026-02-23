# NPCgenGo â€” running and data directory

Quick notes on running the binary and where it looks for `data/`.

## Project structure

Applied structure scope (only):

- `cmd/` â€” executable entrypoints and CLI wiring
- `pkg/` â€” reusable core application logic (`pkg/product/...`)
- `internal/` â€” module-private implementation details (`internal/platform/...`)

No additional optional folders from that setup were added.

Current project folders also include:
- `data/` â€” runtime JSON creation data and local NPC database files
- `web_demo/` â€” web demonstration assets

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
