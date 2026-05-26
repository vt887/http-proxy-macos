# MacNet Gateway (PR-1 Scaffold)

MacNet Gateway is a local macOS network control plane architecture with:

- native SwiftUI desktop app
- Go local daemon (`macnet-gatewayd`)
- optional privileged helper (`macnet-helper`)

PR-1 intentionally ships a clean scaffold with mock integrations (no real Squid/DNS control yet).

## Repository layout

- `app/MacNetGateway`: SwiftUI package with sidebar and core pages
- `daemon`: Go daemon with local API, SQLite init, mock metrics/services/events
- `helper`: privileged helper skeleton
- `docs`: architecture + security + implementation docs
- `packaging`: launchd and install script placeholders

## Local development

### Run daemon

```bash
make daemon-dev
```

Default bind: `127.0.0.1:18080`

### Test daemon

```bash
make daemon-test
```

### Test Swift app module

```bash
cd app/MacNetGateway
swift test
```

### Example API calls

```bash
curl -s http://127.0.0.1:18080/api/health
curl -s http://127.0.0.1:18080/api/dashboard
curl -s http://127.0.0.1:18080/api/services
curl -s http://127.0.0.1:18080/api/live-activity
curl -s http://127.0.0.1:18080/api/settings
curl -s -X PATCH http://127.0.0.1:18080/api/settings -d '{"ui.theme":"dark"}'
```

## PR-2 progress

1. Swift app now includes `DaemonAPIClient` (loopback HTTP calls to daemon endpoints).
2. Dashboard and Services pages now use explicit load states (`idle/loading/loaded/failed`) with retry actions.
3. Daemon persists and serves service statuses from SQLite (`service_status` table) and seeds defaults at startup.
4. Settings screen now loads and saves `ui.theme` through daemon API.
5. Live Activity page now renders daemon events from `/api/live-activity`.
