# Repository instructions for GitHub Copilot

## Project summary

MacNet Gateway is a macOS local network control-plane product split into:
- SwiftUI desktop app (`app/MacNetGateway`)
- Go daemon (`daemon`) exposing local API on loopback
- optional helper skeleton (`helper`)

The daemon is the source of truth for persisted settings/service status (SQLite), while SwiftUI consumes daemon API and renders UI.

## Safety and architecture constraints

- Keep daemon bound to `127.0.0.1` by default.
- Do not expose admin APIs to LAN by default.
- Do not add privileged/root logic to the SwiftUI app.
- Keep app/daemon/helper boundaries explicit.
- Prefer service interfaces and mocks over direct system coupling in early PRs.

## Build and test commands

Run from repository root:

```bash
make format
make daemon-test
cd app/MacNetGateway && swift test
```

Or combined:

```bash
make test
```

For local daemon run:

```bash
make daemon-dev
```

## Key paths

- App entry: `app/MacNetGateway/Sources/MacNetGateway/App/MacNetGatewayApp.swift`
- App API client: `app/MacNetGateway/Sources/MacNetGateway/Services/APIClient.swift`
- App view models: `app/MacNetGateway/Sources/MacNetGateway/ViewModels/`
- Daemon entry: `daemon/cmd/macnet-gatewayd/main.go`
- Daemon HTTP routes: `daemon/internal/api/server.go`
- Daemon DB logic: `daemon/internal/db/db.go`
- Service interfaces: `daemon/internal/services/interfaces.go`

## Change expectations

- Keep changes small and PR-scoped (PR-1 bootstrap, PR-2 app-daemon integration, PR-3 service skeletons).
- Prefer explicit loading/error states in SwiftUI.
- For daemon handlers, use request context (`r.Context()`) for DB calls.
- Add/adjust tests when behavior changes.
