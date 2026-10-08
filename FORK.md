# loveboat/elsbrock-plundrio

Fork of [elsbrock/plundrio](https://github.com/elsbrock/plundrio) carrying one
fix until it lands upstream.

## Why

Upstream flattens put.io's folder tree when listing a transfer's files, so a
pack with per-episode subtitle folders (`Subs/<episode>/2_eng.srt`, ...) maps
every subtitle to the same local path and the whole transfer fails with
`multiple transfer files map to local path`.

The fix keeps each file's path relative to the transfer root
(`fix(api): keep subdirectory paths in transfer file names`). The same bug was
fixed independently in [doodla/plundrio](https://github.com/doodla/plundrio)
(`d7f855a`), which has since diverged from upstream; this fork reapplies the
idea on current upstream instead.

## Building

```bash
docker build --build-arg VERSION=0.11.2-subdirs.1 -t elsbrock-plundrio:0.11.2-subdirs.1 .
```

`Dockerfile` is fork-only; upstream builds with Nix (`flake.nix`).

## Keeping up to date

```bash
git fetch upstream
git merge upstream/main
go test ./...
```

Drop this fork once upstream has the fix.
