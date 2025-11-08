# CLIProxyAPI Fork Notes: Multi-Account Support & Binary Artifact

This fork adds VibeProxy-compatible multi-account auth handling plus a deterministic binary path.

Key behaviors
- Per-account auth files in `~/.cli-proxy-api` named `provider-accountId.json`.
- `type` and `accountId` always set; legacy files (`provider.json`) still read as single accounts.
- Active account chosen via `~/.cli-proxy-api/active-accounts.json` (provider -> id/email/filename hints), with graceful fallback.
- Expired accounts (via `expired`) avoided when valid accounts exist.
- `accountNickname` and unknown fields preserved when updating files.

Implementation highlights
- New `internal/accounts` package:
  - Scans `~/.cli-proxy-api/*.json` (excluding `active-accounts.json`) into `ProviderAccounts`.
  - Infers provider + `accountId` from JSON and filename.
  - Loads `active-accounts.json` and resolves active account per provider with VibeProxy-compatible matching rules.
  - `SaveProviderAccount` writes/merges `provider-accountId.json` preserving metadata.
- Provider auth storages (Codex, Claude, Gemini, Qwen) now save via `SaveProviderAccount` using stable IDs (e.g. account id/email) instead of single global files.
- Request executors are wired to use resolved accounts (and fall back to legacy env/config behavior when needed).

Binary artifact
- Standard build target for automation:
  - `go build -o ./bin/cli-proxy-api ./cmd/...`
- Tools (e.g. VibeProxy) can download the compiled binary from `bin/cli-proxy-api` in this repo.

Upstream merge guidance
- Conflicts are expected primarily under `internal/auth/*`, `internal/runtime/executor/*`, and the new `internal/accounts` package.
- When updating from upstream:
  - Keep `internal/accounts` intact.
  - Ensure any upstream auth changes continue to call `SaveTokenToFile` so per-account files remain the source of truth.
  - Verify `active-accounts.json` resolution and binary path still behave as specified in `CLIProxyAPI_MULTI_ACCOUNT_SPEC.md`.
