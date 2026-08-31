# Implementation Plan: Search Switch (`-s`)

## Overview

Add a `--search` (`-s`) flag to the Recall CLI's root command that runs the same search as the `search` subcommand. Refactor the search-and-print logic into a shared `runSearch(query)` helper in `cmd/search.go`, register the flag in `cmd/root.go`, and branch in `runRecall`. The flag is mutually exclusive with `--edit`, compatible with `--raw`, and registered as a local flag so subcommands reject it.

## Tasks

- [x] 1. Refactor shared search logic
  - [x] 1.1 Extract `runSearch(query string) error` in `cmd/search.go`
    - Move the directory resolution, `search.Search` call, and result-printing out of `searchCmd.RunE` into `runSearch`
    - Have `searchCmd.RunE` keep its no-argument guard, then call `runSearch(args[0])`
    - Preserve exact output format (`filename:line:content`, `----------` separators) and exit behavior
    - _Requirements: 2.2, 2.3, 2.4, 6.1, 6.2_

- [x] 2. Register the search flag
  - [x] 2.1 Declare `searchFlag` and register in `init()`
    - Add `var searchFlag bool` at package level in `cmd/root.go`
    - Add `rootCmd.Flags().BoolVarP(&searchFlag, "search", "s", false, "search all recall files for the given query")` in `init()`
    - _Requirements: 1.1, 1.2, 1.3, 1.4_

- [x] 3. Wire `-s` into `runRecall`
  - [x] 3.1 Add exclusivity check and search branch
    - After the `len(args) == 0` help check, add `if searchFlag && editFlag { print error; os.Exit(1) }`
    - Then `if searchFlag { return runSearch(args[0]) }`
    - Leave the existing `--raw`/`--edit` exclusivity, edit delegation, and display path unchanged
    - _Requirements: 2.1, 3.1, 4.1, 4.2_

- [x] 4. Checkpoint - Build and manual smoke test
  - Run `make build`; verify `recall -s "..."` matches `recall search "..."`
  - _Requirements: 2.2_

- [x] 5. Write integration tests in `cmd/root_test.go`
  - [x] 5.1 Add subprocess tests using existing helpers
    - `-s "Hello"` output equals `search "Hello"` output
    - `--search "Hello"` output equals `search "Hello"` output
    - `-s "zzznomatch"` → no output, exit 0
    - `-s` with no query → help, exit 0
    - `-s --edit hello` → exclusivity error, exit 1
    - `-s -r "Hello"` → accepted, exit 0
    - `--search`/`-s` rejected on edit/list/search/init subcommands
    - _Requirements: 1.1, 1.2, 1.3, 2.2, 2.4, 3.1, 4.1, 4.2, 5.1–5.6_

- [x] 6. Update README
  - [x] 6.1 Add `recall -s "docker compose"` example to the Search section, noting equivalence to `recall search`
    - _Requirements: 1.3, 2.2_

- [x] 7. Final checkpoint - Ensure all tests pass
  - Run `make build` and `go test ./...`

## Notes

- Each task references specific requirements for traceability.
- All implementation changes are confined to `cmd/search.go` and `cmd/root.go`; tests go in `cmd/root_test.go`.
- The search logic itself (`internal/search`) is reused unchanged.

## Task Dependency Graph

```json
{
  "waves": [
    { "id": 0, "tasks": ["1.1", "2.1"] },
    { "id": 1, "tasks": ["3.1"] },
    { "id": 2, "tasks": ["5.1", "6.1"] }
  ]
}
```
