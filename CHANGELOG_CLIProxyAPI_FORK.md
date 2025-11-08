# CLIProxyAPI Fork Changelog

## Unreleased

- Add `internal/accounts` package to implement VibeProxy-compatible multi-account storage:
  - Load per-account auth JSON from `~/.cli-proxy-api/provider-accountId.json`.
  - Infer `type`/`accountId` from JSON and filenames, skipping malformed files.
  - Load `active-accounts.json` and resolve the active account per provider with fallback.
  - Provide `SaveProviderAccount` helper that merges updates into existing files without dropping metadata.
- Update auth token storage implementations to write per-account files via `SaveProviderAccount`:
  - Codex: persist by `AccountID`.
  - Claude: persist by email as `accountId`.
  - Gemini: persist by email as `accountId`.
  - Qwen: persist to a stable provider-specific account id.
- Ensure expired accounts (via `expired`) are avoided when valid options exist; preserve `accountNickname` and unknown fields.
- Introduce a deterministic binary build path for automation:
  - `go build -o ./bin/cli-proxy-api ./cmd/...` for consumers (e.g. VibeProxy) to download.
