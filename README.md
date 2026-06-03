# Leonidas C2 Framework

A modified fork of the [Sliver](https://github.com/BishopFox/sliver) adversary emulation / red team framework.

Module: `github.com/leonidas-c2/leonidas`

## About

Leonidas keeps Sliver's core architecture (mTLS, WireGuard, HTTP(S), DNS C2, dynamic implants, multiplayer, etc.) with project-specific changes such as ICMP C2, a gRPC-Web web UI, and related server/client updates.

Upstream Sliver is [GPLv3](https://www.gnu.org/licenses/gpl-3.0.html). This fork remains under the same license; see `LICENSE` and upstream attribution.

## Build

```bash
make          # server + client (leonidas-server, leonidas-client)
make client   # client only
./go-tests.sh --unit-only   # faster test run (skips generate + c2 e2e)
```

See `AGENTS.md` for development commands.

## Documentation

- [Sliver docs](https://sliver.sh/docs) — most workflows still apply
- [BishopFox/sliver](https://github.com/BishopFox/sliver) — upstream project

## Disclaimer

Use only on systems and networks you are authorized to test.
