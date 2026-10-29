# Etan's CLIProxyAPI carry

Base: upstream `v8.0.45`, matching the Homebrew hub version when this carry was created.
Branch: https://github.com/EtanHey/CLIProxyAPI/tree/etan/carry

Upstream fix: https://github.com/router-for-me/CLIProxyAPI/pull/6700
(`cyhhao:fix/error-ratelimit-headers`). No replacement upstream issue or PR.

## Carried commits

The six upstream commits were replayed in order onto the release tag. The author
and individual commits are retained; conflicts preserve the release's newer
executor fields and streaming tests.

| Upstream commit | Carry commit | Change |
|---|---|---|
| d14a4621 | 63676173 | Carry retry headers on Claude errors |
| 9405a7b7 | ef09f8d1 | Preserve headers through fast error wrappers |
| 2c091e90 | f2b3613b | Find headers through execution wrappers |
| 4755d85d | df572be7 | Filter error headers |
| b59f0634 | 01139bc5 | Centralize filtering |
| 860a2f44 | 9165018c | Centralize retry timing precedence |

Local fix `fcce7f91` following these commits covers the direct-response and
pre-stream rejection paths. Its regression tests enforce the rate-header
allowlist and exclude authentication and cookie headers. See `git log v8.0.45..etan/carry`.

The local allowlist is `Retry-After` and `Anthropic-Ratelimit-*`. Claude 429
responses expose this recovery metadata even when generic header passthrough is
off, so the existing service config can stay unchanged. Fast direct errors retain
only this subset and their Content-Type. CPA's trusted local Retry-After still
takes precedence. Upstream reset headers describe the last credential tried;
they do not establish a pooled-account recovery deadline.

Cache TTL and affinity policy are unchanged by this carry.

## Verification

Use the Go toolchain pinned by go.mod. The local fix was preceded by failing
regressions for direct/normal errors, streaming/nonstreaming rejection,
passthrough on/off, and the header allowlist. Run the full suite and compile:

```sh
go test ./...
go build -o cli-proxy-api ./cmd/server
```

The carry regressions and the full suites of its four changed packages pass on
the M1. The local full-suite run reports failures in catalog/registry config,
the byte-write invariant, and Pion media relay tests. Their receipts and the
untouched-tag comparison are kept in the local C6 report.

## Return to upstream

When upstream merges and Homebrew ships it, rebase to upstream or drop the carry
and return to brew services. Check weekly:

```sh
gh pr view 6700 --repo router-for-me/CLIProxyAPI --json state,mergedAt
brew info cliproxyapi
```

Before dropping the carry, rerun the synthetic error-header probe on the new
Homebrew binary, including nonstream and stream-start 429 responses. A merged PR
alone does not prove that Homebrew contains it or that every carried path works.

## Runtime and rollback

The carry uses a separate user LaunchAgent, `com.tokenlayer.cliproxyapi`, with
the existing `/opt/homebrew/etc/cliproxyapi.conf`. Homebrew's binary and plist
are retained. Runtime evidence and exact installed commit belong in the local
C6 report, rather than this public tracking file.

Rollback on the M1:

```sh
launchctl bootout "gui/$(id -u)/com.tokenlayer.cliproxyapi"
launchctl disable "gui/$(id -u)/com.tokenlayer.cliproxyapi"
brew services start cliproxyapi
```

The disable step keeps the carry from starting again at the next login.
Homebrew's start command enables its own label before bootstrapping it.

Do not leave both services listening on port 8317. To use a new carry build,
stop the carry agent before switching its binary, then enable and bootstrap it again.
