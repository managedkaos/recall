# Implementation Plan: Metadata Switch (`-m` / `--metadata`)

## Overview

Add a `--metadata` (`-m`) action flag and a `--json` (`-j`) modifier to the
Recall CLI's root command. For each named file that exists, report its absolute
path, filesystem metadata (size, modification time, and platform-available
creation/change times), and front-matter tags — as human-readable text by
default or as a JSON array with `--json`. Implement the handler in a new
`cmd/metadata.go`, isolate platform-specific time resolution behind build-tagged
helpers, register the flags in `cmd/root.go`, and wire a dispatch branch into
`runRecall`. `--metadata` is an action flag subject to the existing
one-action-flag rule; `--json` is valid only alongside `--metadata`.

## Tasks

- [ ] 1. Platform time helpers
  - [ ] 1.1 Add the shared signature and per-OS implementations
    - Create `cmd/metadata_time_darwin.go`, `_linux.go`, `_windows.go`, `_other.go` with build tags
    - Each exposes `platformTimes(info os.FileInfo) (created, changed *time.Time)`
    - darwin: `Birthtimespec` → created, `Ctimespec` → changed
    - linux: `Ctim` → changed, created nil
    - windows: `CreationTime` → created, changed nil
    - other: return `nil, nil`
    - Guard every `Sys()` type assertion with comma-ok; never panic on nil/unexpected
    - _Requirements: 2.4_

- [ ] 2. Metadata collection and rendering
  - [ ] 2.1 Define `FileMetadata` and `collectMetadata` in `cmd/metadata.go`
    - Struct with `Name`, `Path`, `SizeBytes`, `Modified`, `Created *time.Time`, `Changed *time.Time`, `Tags` and JSON tags
    - `collectMetadata` combines `filepath.Abs`, `os.Stat`, `platformTimes`, `storage.Read`, `frontmatter.Parse`
    - Normalize a nil tag slice to an empty (non-nil) slice
    - _Requirements: 2.2, 2.3, 2.4, 2.5_
  - [ ] 2.2 Implement `renderText` and `renderJSON`
    - `renderText`: labeled block per file, omit Created/Changed when nil, "(none)" for empty tags, blank line between files
    - `renderJSON`: normalize nil slice to `[]`, `json.Encoder` with two-space indent, RFC 3339 timestamps
    - _Requirements: 2.2, 2.3, 4.1, 4.2, 4.3, 4.4, 4.5_

- [ ] 3. Metadata handler
  - [ ] 3.1 Implement `runMetadata(names []string, asJSON bool) error`
    - Resolve/ensure recall dir (error → stderr, exit 1)
    - For each name: `storage.Exists` false → stderr error, continue; else `collectMetadata`
    - Dispatch to `renderJSON` or `renderText`
    - Exit 1 only when no file was collected; otherwise return nil
    - _Requirements: 2.1, 2.6, 3.1, 3.2, 3.3, 3.4, 7.1, 7.2_

- [ ] 4. Register flags and wire dispatch in `cmd/root.go`
  - [ ] 4.1 Declare `metadataFlag` and `jsonFlag`; register in `init()`
    - Add `metadataFlag` to the action-flag group; `jsonFlag` as modifier
    - `f.BoolVarP(&metadataFlag, "metadata", "m", false, ...)`
    - `f.BoolVarP(&jsonFlag, "json", "j", false, ...)`
    - _Requirements: 1.1, 1.2, 1.3, 1.4, 1.5, 1.6, 6.1_
  - [ ] 4.2 Enforce exclusivity, `--json` dependency, and add dispatch branch
    - Add `metadataFlag` to the action-count slice (one-action-flag rule)
    - After the `--init-path` guard, add `if jsonFlag && !metadataFlag { error; exit 1 }`
    - Add `case metadataFlag:` → no args ⇒ help; else `runMetadata(args, jsonFlag)`
    - _Requirements: 3.5, 4.6, 5.1, 5.2, 5.3_

- [ ] 5. Checkpoint - Build and manual smoke test
  - Run `make build`
  - Verify `recall -m <file>`, `recall -m -j <file>`, missing-file, and no-arg behaviors by hand
  - _Requirements: 2.2, 3.5, 4.1_

- [ ] 6. Unit tests (`cmd/metadata_test.go`)
  - [ ] 6.1 Test `collectMetadata`, `renderJSON`, and `platformTimes`
    - `collectMetadata` on a tagged temp file: path, size, modtime, tags
    - tagless file ⇒ empty non-nil tag slice
    - `renderJSON([])` writes exactly `[]\n`; populated slice round-trips via `json.Unmarshal`
    - `platformTimes` on a temp file: no panic; non-nil returns are plausible
    - _Requirements: 2.4, 2.5, 4.2, 4.5_

- [ ] 7. Integration tests (`cmd/root_test.go`)
  - [ ] 7.1 Add subprocess tests using existing helpers
    - `--metadata`/`-m` and `--json`/`-j` shown in `--help`
    - `-m docker` → path/size/modtime/tags, exit 0
    - `-m -j docker` → parseable JSON array with expected fields
    - `-m nope` → stderr error, exit 1
    - `-m docker nope` → docker reported, stderr error, exit 0
    - `-m -j nope` → stdout `[]`, exit 1
    - `-m` no args → help, exit 0
    - `-j docker` without `-m` → error, exit 1
    - `-m -l` → one-action-flag error, exit 1
    - `-m`/`-j` rejected on subcommands
    - directory contents unchanged after `-m` run
    - _Requirements: 1.5, 2.2–2.6, 3.1–3.5, 4.1, 4.2, 4.4, 4.5, 4.6, 5.2, 6.1, 6.2, 7.1_

- [ ] 8. Update documentation
  - [ ] 8.1 Add a "Report metadata" section to README and a row to the action table
    - Document `recall -m <file>`, multi-file usage, and `-m -j` JSON output
    - Note platform differences for creation/change times
    - _Requirements: 1.5_
  - [ ] 8.2 Update ROADMAP if the feature was listed there
    - _Requirements: n/a_

- [ ] 9. Final checkpoint - Ensure all tests pass
  - Run `make build` and `go test ./...`

## Notes

- Each task references specific requirements for traceability.
- Implementation is confined to new files `cmd/metadata.go` and
  `cmd/metadata_time_*.go`, plus edits to `cmd/root.go`; tests go in
  `cmd/metadata_test.go` and `cmd/root_test.go`.
- Existing packages (`config`, `storage`, `frontmatter`) are reused unchanged.
- Cross-platform time fields are the main risk; keep all OS-specific code behind
  build tags with defensive type assertions.

## Task Dependency Graph

```json
{
  "waves": [
    { "id": 0, "tasks": ["1.1", "2.1"] },
    { "id": 1, "tasks": ["2.2", "3.1"] },
    { "id": 2, "tasks": ["4.1", "4.2"] },
    { "id": 3, "tasks": ["6.1", "7.1", "8.1", "8.2"] }
  ]
}
```
