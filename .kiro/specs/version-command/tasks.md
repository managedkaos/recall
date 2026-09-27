# Implementation Checklist: Version Metadata

- [x] Inject a single `cmd.Version` from Make or GoReleaser.
- [x] Derive local versions from Git with explicit override and unknown fallback.
- [x] Normalize version display while preserving development/release suffixes.
- [x] Preserve `--version` / `-v` and all other build metadata.
- [x] Fetch tags/history in CI Make builds and align release checkout refs.
- [x] Remove the separate version file and post-release synchronization job.
- [x] Document local, plain Go, snapshot, and release version behavior.
- [x] Validate unit/CLI tests, Make Git fixtures, cross-platform builds,
      GoReleaser configuration and snapshot, and workflow syntax.
