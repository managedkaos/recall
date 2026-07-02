# Implementation Plan: Raw Output Flag

## Overview

Add a `--raw` (`-r`) flag to the Recall CLI's root command in `cmd/root.go`. The flag bypasses the Glamour renderer and writes the front-matter-stripped body directly to stdout. It is mutually exclusive with `--edit` and registered as a local flag so subcommands are unaffected.

## Tasks

- [x] 1. Add raw flag variable and registration
  - [x] 1.1 Declare `rawFlag` variable and register the flag in `init()`
    - Add `var rawFlag bool` at package level in `cmd/root.go`
    - Add `rootCmd.Flags().BoolVarP(&rawFlag, "raw", "r", false, "output unformatted markdown without ANSI styling")` in `init()`
    - _Requirements: 1.1, 1.2, 1.3, 4.5_

- [x] 2. Implement mutual-exclusivity check and raw output branch
  - [x] 2.1 Add mutual-exclusivity guard in `runRecall()`
    - After extracting the filename argument, check if both `rawFlag` and `editFlag` are true
    - If both are set, print `"recall: --raw and --edit flags cannot be used together"` to stderr and `os.Exit(1)`
    - _Requirements: 1.4_

  - [x] 2.2 Add raw output branch after `frontmatter.Parse()`
    - After the existing `_, body := frontmatter.Parse(content)` line, add a conditional: if `rawFlag` is true, call `os.Stdout.Write(body)` and return nil
    - The existing `renderer.Render(body)` path remains in the else branch
    - _Requirements: 2.1, 2.2, 2.3, 2.4_

- [x] 3. Checkpoint - Verify flag registration and basic behavior
  - Ensure all tests pass, ask the user if questions arise.

- [x] 4. Write tests for raw output
  - [x] 4.1 Write unit tests for raw flag behavior in `cmd/root_test.go`
    - Test that `--raw` with a valid file produces the body content on stdout and exits 0
    - Test that `-r` shorthand works identically
    - Test that `--raw` and `--edit` together produce an error exit
    - Test that `--raw` without a filename argument shows help and exits 0
    - Test that `--raw` with a non-existent file exits 1 silently
    - Test that `--raw`/`-r` is rejected on edit, list, search, and init subcommands
    - _Requirements: 1.1, 1.2, 1.4, 1.5, 2.4, 3.1, 3.2, 4.1–4.4, 4.6_

  - [ ]* 4.2 Write property test: Raw output equals parsed body
    - **Property 1: Raw output equals parsed body**
    - **Validates: Requirements 2.1, 2.3**
    - Use `pgregory.net/rapid` to generate arbitrary byte slices (with and without front-matter tags line)
    - Assert that passing content through `frontmatter.Parse` and writing the body to a buffer produces output identical to the body returned by `frontmatter.Parse`

  - [ ]* 4.3 Write property test: Raw output contains no ANSI escape sequences
    - **Property 2: Raw output contains no ANSI escape sequences**
    - **Validates: Requirements 2.2**
    - Use `pgregory.net/rapid` to generate arbitrary byte slices (with and without front-matter tags line)
    - Assert that the raw output path never introduces bytes matching the ANSI escape pattern `\x1b[`

- [x] 5. Final checkpoint - Ensure all tests pass
  - Ensure all tests pass, ask the user if questions arise.

## Notes

- Tasks marked with `*` are optional and can be skipped for faster MVP
- Each task references specific requirements for traceability
- Checkpoints ensure incremental validation
- Property tests validate universal correctness properties using `pgregory.net/rapid`
- Unit tests validate specific examples and edge cases
- All changes are confined to `cmd/root.go` (implementation) and a new `cmd/root_test.go` (tests)

## Task Dependency Graph

```json
{
  "waves": [
    { "id": 0, "tasks": ["1.1"] },
    { "id": 1, "tasks": ["2.1", "2.2"] },
    { "id": 2, "tasks": ["4.1", "4.2", "4.3"] }
  ]
}
```
