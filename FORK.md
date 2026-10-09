# loveboat/elsbrock-plundrio

Fork of [elsbrock/plundrio](https://github.com/elsbrock/plundrio) carrying two
fixes until they land upstream.

## Why

**Subtitle folders.** Upstream flattens put.io's folder tree when listing a
transfer's files, so a pack with per-episode subtitle folders
(`Subs/<episode>/2_eng.srt`, ...) maps every subtitle to the same local path
and the whole transfer fails with `multiple transfer files map to local path`.

The fix keeps each file's path relative to the transfer root
(`fix(api): keep subdirectory paths in transfer file names`). The same bug was
fixed independently in [doodla/plundrio](https://github.com/doodla/plundrio)
(`d7f855a`), which has since diverged from upstream; this fork reapplies the
idea on current upstream instead.

**Stale put.io errors.** put.io keeps a transfer's `error_message` after it
recovers: a transfer queued while the account was full still says `You need
6.9 G free space to start this transfer.` once it has completed. Upstream
passes that to Sonarr/Radarr as a Transmission error, so they show the
finished download as a warning and never import it. The fix drops put.io's
error once the transfer is `COMPLETED` or `SEEDING`; local errors are still
reported (`fix(server): ignore stale put.io error on finished transfers`).

## Building

```bash
docker build --build-arg VERSION=0.11.2-subdirs.2 -t elsbrock-plundrio:0.11.2-subdirs.2 .
```

`Dockerfile` is fork-only; upstream builds with Nix (`flake.nix`).

## Keeping up to date

```bash
git fetch upstream
git merge upstream/main
go test ./...
```

Drop this fork once upstream has both fixes.
