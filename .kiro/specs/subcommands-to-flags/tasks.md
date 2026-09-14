# Implementation Plan: Subcommands to Flags

## Overview

Remove all Cobra subcommands (`edit`, `list`, `search`, `init`, `version`,
built-in `completion`) and reimplement each as a root-level flag. Short switch
= first letter of the former subcommand (`-e -s -l -i -v -c`); long switch =
the spelled-out name (`--edit --search --list --init --version --completion`).
`--completion` takes a shell value; the rest are boolean except the existing
`--raw` and the promoted `--tag`. Centralize action dispatch and
mutual-exclusivity in `runRecall`. Remove reserved-name enforcement.

## Tasks

- [ ] 1. Convert helpers to be subcommand-independent
  - [ ] 1.1 `cmd/list.go`: extract `runList(tag string) error` from `listCmd.RunE`; remove `listCmd`, its `--tag` flag registration, and its `AddCommand`
    - _Requirements: 4.2, 4.3, 1.2_
  - [ ] 1.2 `cmd/init.go`: extract `runInit(path string) error` (empty path → default/`RECALL_DIR`; non-empty → that dir); remove `initCmd` and its `AddCommand`
    - _Requirements: 5.2, 5.3, 1.4_
  - [ ] 1.3 `cmd/version.go`: remove `versionCmd` and its `AddCommand`; keep `runVersion`
    - _Requirements: 6.2, 1.5_
  - [ ] 1.4 `cmd/search.go`: remove `searchCmd` and its `AddCommand`; keep `runSearch`
    - _Requirements: 3.2, 1.3_
  - [ ] 1.5 `cmd/edit.go`: remove `editCmd` and its `AddCommand`; delete the `IsReservedName` check inside `runEdit`
    - _Requirements: 2.2, 1.1, 1.8_

- [ ] 2. Add completion generator
  - [ ] 2.1 Create `cmd/completion.go` with `runCompletion(cmd, shell)` covering bash/zsh/fish/powershell and an error for unsupported/empty shells
    - _Requirements: 7.1, 7.2, 7.3, 7.4_

- [ ] 3. Rework the root command
  - [ ] 3.1 In `cmd/root.go`, declare all flag vars (`listFlag`, `initFlag`, `versionFlag`, `completionFlag`, `tagFlag`) alongside existing ones
    - _Requirements: 8.1, 8.2_
  - [ ] 3.2 Register all flags in `init()`; set `rootCmd.CompletionOptions.DisableDefaultCmd = true`
    - _Requirements: 1.6, 8.1, 8.2, 8.3_
  - [ ] 3.3 Extract the current default render body into `renderFile(name string) error`
    - _Requirements: 9.3_
  - [ ] 3.4 Rewrite `runRecall` with the single action-count check and the action `switch`
    - _Requirements: 9.1, 9.2, 9.3, 9.4, 2.3, 3.3_
  - [ ] 3.5 Remove `reservedNames` and `IsReservedName`
    - _Requirements: 1.8_

- [ ] 4. Checkpoint - build and smoke test
  - Run `make build`
  - Verify: `recall --help` shows flags only; `recall list` treats `list` as a filename; `recall -s Hello`, `recall -l`, `recall -v`, `recall -i`, `recall -c bash` behave
  - _Requirements: 1.7, 10.1_

- [ ] 5. Tests
  - [ ] 5.1 Update `cmd/root_test.go`: remove subcommand-presence and reserved-name assertions
    - _Requirements: 1.1–1.8_
  - [ ] 5.2 Add flag-based integration tests per the design's test matrix (edit/search/list/init/version/completion, conflict, short forms)
    - _Requirements: 2.x, 3.x, 4.x, 5.x, 6.x, 7.x, 8.1, 9.1, 10.1_
  - [ ] 5.3 Update/remove `cmd/version_test.go` if it invokes the `version` subcommand
    - _Requirements: 6.2_

- [ ] 6. Documentation
  - [ ] 6.1 Update `README.md`: replace subcommand usage with flag usage; document `-c <shell>`, `-i`, `-l`/`--tag`, `-v`; note `init` path now via `RECALL_DIR` only
    - _Requirements: 10.2, 5.3_
  - [ ] 6.2 Update `ROADMAP.md` if the `-s` item (now generalized) needs reframing
    - _Requirements: 10.2_

- [ ] 7. Final checkpoint
  - Run `make build` and `go test ./...`; all pass

## Resolved Decisions

1. **`init` custom directory** — via `--init-path <dir>` flag (approved).
2. **`completion` as `-c <shell>` string flag** (approved).
3. **`reservedNames` removed entirely** — no filename is rejected (approved).

## Task Dependency Graph

```json
{
  "waves": [
    { "id": 0, "tasks": ["1.1", "1.2", "1.3", "1.4", "1.5", "2.1"] },
    { "id": 1, "tasks": ["3.1", "3.2", "3.3", "3.4", "3.5"] },
    { "id": 2, "tasks": ["4"] },
    { "id": 3, "tasks": ["5.1", "5.2", "5.3", "6.1", "6.2"] },
    { "id": 4, "tasks": ["7"] }
  ]
}
```
