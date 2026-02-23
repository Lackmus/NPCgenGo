# Wails Frontend Assets

This folder contains desktop UI assets used by the Wails host in `cmd/npcgen-wails`.

## Build prerequisites

- Go installed
- Wails CLI installed:
  - `go install github.com/wailsapp/wails/v2/cmd/wails@v2.10.2`
- Ensure `GOPATH/bin` is in your PATH (Windows example: `C:\Users\<you>\go\bin`)

## Build commands

Run from the project root:

- Wails CLI build (recommended):
  - `Push-Location cmd/npcgen-wails; wails build -skipbindings -s; Pop-Location`

- Manual build (from Wails manual builds guide):
  - `Push-Location cmd/npcgen-wails; go build -tags "desktop,production" -ldflags "-w -s -H windowsgui" -o build/bin/npcgen-wails-manual.exe; Pop-Location`

## Notes

- Wails config is in `cmd/npcgen-wails/wails.json`.
- Frontend assets are static in `ui/wails/dist` and embedded by `ui/wails/assets.go`.
- Data directory resolution for `cmd/npcgen-wails`:
  - `--data-dir <path>` (highest precedence)
  - `NPCGEN_DATA` environment variable
  - auto-discovery by searching upward for `data/creation_data/factiondata` from CWD, then executable dir
  - final fallback: relative `data`
- Building from subfolders is supported, e.g. `cd cmd/npcgen-wails; wails build` will still find project-root `data/`.
- If `wails` is still not found, either restart the terminal after PATH update or run via absolute path:
  - `C:\Users\<you>\go\bin\wails.exe`

## Troubleshooting

If you see:

`Wails applications will not build without the correct build tags.`

You launched it with plain `go build`/`go run` (missing tags). Use one of:

- `cd cmd/npcgen-wails; wails build -clean`
- `cd cmd/npcgen-wails; go build -tags "desktop,production" -ldflags "-w -s -H windowsgui" -o build/bin/npcgen-wails-manual.exe`
- For local run via Go toolchain (without `wails dev`): `cd cmd/npcgen-wails; go run -tags "desktop,production" .`
- Use `desktop,dev` only through `wails dev`.
