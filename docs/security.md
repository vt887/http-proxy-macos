# Security Decisions (PR-1)

## Enforced defaults

- Daemon binds to loopback only: `127.0.0.1`.
- No LAN admin API and no remote management in PR-1.
- No privileged helper execution path in PR-1.
- No SSL MITM, no certificate installation logic.

## Config safety

- Settings writes are explicit (`PATCH /api/settings`).
- Database path is under user-scoped dev directory in PR-1 to avoid writing into repository files.

## Future hardening (PR-2+)

- Add explicit auth + opt-in before any non-loopback API exposure.
- Add audit logging for privileged workflows once helper actions exist.
- Add config backup + validation gates before applying Squid/DNS changes.
