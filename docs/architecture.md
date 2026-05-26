# MacNet Gateway Architecture (PR-1)

## Boundaries

- **SwiftUI app**: desktop UX only (sidebar, dashboard, services, placeholders).
- **Go daemon (`macnet-gatewayd`)**: local control plane API, SQLite init, mock metrics/events/services.
- **Helper (`macnet-helper`)**: placeholder only; no privileged operations yet.

## Runtime communication

- PR-1 uses local HTTP API on `127.0.0.1:18080`.
- LAN binding is intentionally disabled by default.

## Implemented API

- `GET /api/health`
- `GET /api/dashboard`
- `GET /api/services`
- `GET /api/live-activity`
- `GET /api/settings`
- `PATCH /api/settings`

## Persistence

- SQLite initialized on startup (dev default path in `~/.macnet-gateway-dev/db/app.sqlite`).
- PR-1 tables: `settings`, `service_status`, `audit_log`.

## PR-2 focus

- Replace Swift mock client with real daemon API client.
- Add loading/error states across pages.
- Persist and expose service statuses and settings via daemon as first-class models.
