This tiny demo shows how to load your `data/creation_data` JSON files and store generated NPCs in the browser's IndexedDB using the `idb` library.

How to run:
1. Serve the demo from your Go backend (recommended).

This demo now uses your Go backend API for NPC storage. Make sure your Go server serves static files from `ui/web` (mounted at `/ui`) and exposes the following endpoints:

- `GET  /api/npcs`        â€” return list of stored NPCs (JSON array)
- `POST /api/npcs`        â€” create NPC (accepts NPC JSON body)
- `DELETE /api/npcs/:id`  â€” delete NPC by id

Example: run your Go app (it listens on `:8080` in the provided server snippet), then open:

http://localhost:8080/ui/

Notes:
- The demo still fetches creation JSON from `data/creation_data/...` (served statically by the Go server).
- If your server uses a different path or port, update `API_BASE` in `ui/web/app.js`.

If you want, I can add a `server.go` file that mounts `ui/web/` and implements the `/api/npcs` handlers using your existing `NPCService` so it works out-of-the-box.
