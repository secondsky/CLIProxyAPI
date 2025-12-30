# CLIProxyAPI Fork Changelog

## Unreleased

- Restore VibeProxy-compatible multi-account support:
  - Per-account auth JSON in `~/.cli-proxy-api/provider-accountId.json`; `active-accounts.json` selects the active account per provider with graceful fallback.
  - Auth storages (Codex/Claude/Gemini/Qwen) persist via `SaveProviderAccount`, keeping metadata intact and skipping expired choices when possible.
  - Watcher synthesizer now honors multi-account files first (default auth dir), avoiding duplicates and still handling legacy single-file auths.
- Deterministic binary build target for automation: `go build -o ./bin/cli-proxy-api ./cmd/...`.
