# Implementation Plan: Recall CLI

## Overview

Implement a Go CLI tool called "Recall" that stores, retrieves, edits, lists, and searches markdown-formatted reference files from the command line. The implementation follows a bottom-up approach: core internal packages first, then CLI commands, and finally cross-compilation and integration wiring.

## Tasks

- [x] 1. Set up project structure and dependencies
  - [x] 1.1 Initialize Go module and install dependencies
    - Run `go mod init` for the recall module
    - Add dependencies: `github.com/spf13/cobra`, `github.com/charmbracelet/glamour`, `pgregory.net/rapid`
    - Create directory structure: `cmd/`, `internal/config/`, `internal/frontmatter/`, `internal/renderer/`, `internal/search/`, `internal/storage/`
    - Create `main.go` entry point that calls `cmd.Execute()`
    - _Requirements: 6.1_

- [x] 2. Implement configuration package
  - [x] 2.1 Implement `internal/config/config.go`
    - Implement `RecallDir()` function that checks `RECALL_DIR` env var first, then falls back to `$HOME/.recall`
    - Implement `EnsureDir(dir string)` function that creates the directory and intermediate directories if they don't exist
    - Return appropriate errors if home directory cannot be determined or directory cannot be created
    - _Requirements: 5.1, 5.2, 5.3, 5.4, 6.2, 6.3_

  - [ ]* 2.2 Write property test for RECALL_DIR configuration
    - **Property 7: RECALL_DIR configuration is respected**
    - **Validates: Requirements 5.1**

  - [ ]* 2.3 Write property test for auto-creation of recall directory
    - **Property 8: Auto-creation of recall directory**
    - **Validates: Requirements 5.3**

  - [ ]* 2.4 Write property test for native path separators
    - **Property 9: Path construction uses native separators**
    - **Validates: Requirements 6.2, 6.4**

- [x] 3. Implement front-matter parser
  - [x] 3.1 Implement `internal/frontmatter/frontmatter.go`
    - Implement `Parse(content []byte) (tags []string, body []byte)` function
    - Implement `ParseTagLine(line string) []string` function
    - Handle edge cases: consecutive commas, trailing commas, whitespace-only values, empty files
    - The `tags:` prefix match must be case-sensitive
    - _Requirements: 8.1, 1.4, 8.3_

  - [ ]* 3.2 Write property test for front-matter exclusion
    - **Property 3: Front-matter exclusion**
    - **Validates: Requirements 1.4, 8.3**

  - [ ]* 3.3 Write property test for tag parsing correctness
    - **Property 11: Tag parsing correctness**
    - **Validates: Requirements 8.1**

- [x] 4. Implement storage package
  - [x] 4.1 Implement `internal/storage/storage.go`
    - Implement `List(dir string) ([]string, error)` — returns sorted non-hidden regular filenames
    - Implement `Read(dir string, name string) ([]byte, error)` — reads file content
    - Implement `FilePath(dir string, name string) string` — constructs OS-native path
    - Implement `Exists(dir string, name string) bool` — checks file existence
    - Implement `Create(dir string, name string) error` — creates empty file
    - _Requirements: 3.1, 1.1, 6.4_

  - [ ]* 4.2 Write property test for list output sorted and filtered
    - **Property 4: List output is sorted and filtered**
    - **Validates: Requirements 3.1**

- [x] 5. Implement search engine
  - [x] 5.1 Implement `internal/search/search.go`
    - Implement `SearchContent(filename string, content []byte, query string) []Result` — case-insensitive substring match within a single file
    - Implement `Search(dir string, query string) ([]FileResults, error)` — scans all files in directory
    - Results must be grouped by file with matches in ascending line-number order
    - _Requirements: 4.1, 4.2, 4.3_

  - [ ]* 5.2 Write property test for search completeness and soundness
    - **Property 5: Search completeness and soundness**
    - **Validates: Requirements 4.1**

  - [ ]* 5.3 Write property test for search output format and ordering
    - **Property 6: Search output format and ordering**
    - **Validates: Requirements 4.2, 4.3, 4.4**

- [x] 6. Implement renderer package
  - [x] 6.1 Implement `internal/renderer/renderer.go`
    - Implement `Render(content []byte) (string, error)` function
    - Use Glamour with "auto" style for terminal-aware rendering
    - Handle edge cases: empty content, invalid markdown
    - _Requirements: 1.3_

- [x] 7. Checkpoint - Core packages complete
  - Ensure all tests pass, ask the user if questions arise.

- [x] 8. Implement CLI commands
  - [x] 8.1 Implement root command (`cmd/root.go`)
    - Set up Cobra root command with application description
    - Implement `recall <filename>` handler: resolve directory, read file, strip front-matter, render, print to stdout
    - Handle non-existent file: exit silently with non-zero code (no stdout/stderr output)
    - Add `-e` flag as shorthand for edit subcommand
    - Wire `cmd.Execute()` for `main.go` to call
    - _Requirements: 1.1, 1.2, 1.3, 1.4, 7.1_

  - [x] 8.2 Implement edit command (`cmd/edit.go`)
    - Register "edit" subcommand with Cobra
    - Resolve recall directory, create file if it doesn't exist
    - Read `$EDITOR` env var; error to stderr if unset
    - Launch editor as subprocess with full file path, wait for exit
    - Propagate editor's exit code
    - _Requirements: 2.1, 2.2, 2.3, 2.4_

  - [x] 8.3 Implement list command (`cmd/list.go`)
    - Register "list" subcommand with Cobra
    - Add `--tag` flag for filtering by tag
    - Resolve recall directory, list files (sorted, no hidden, no dirs)
    - If `--tag` provided: read each file, parse front-matter, filter by case-insensitive tag match
    - Print one filename per line; no output if empty
    - _Requirements: 3.1, 3.2, 3.3, 8.2, 8.4_

  - [x] 8.4 Implement search command (`cmd/search.go`)
    - Register "search" subcommand with Cobra
    - Require search query argument; error to stderr if missing
    - Resolve recall directory, run search engine
    - Format output: `filename:linenum:line_content`, groups separated by `----------`
    - No output and zero exit if no matches
    - _Requirements: 4.1, 4.2, 4.3, 4.4, 4.5, 4.6_

  - [ ]* 8.5 Write property test for recall retrieves stored content
    - **Property 1: Recall retrieves stored content**
    - **Validates: Requirements 1.1**

  - [ ]* 8.6 Write property test for non-existent file silent failure
    - **Property 2: Non-existent file produces silent failure**
    - **Validates: Requirements 1.2**

  - [ ]* 8.7 Write property test for unrecognized input error
    - **Property 10: Unrecognized input produces descriptive error**
    - **Validates: Requirements 7.2**

  - [ ]* 8.8 Write property test for tag filtering
    - **Property 12: Tag filtering returns exactly matching files**
    - **Validates: Requirements 8.2**

- [x] 9. Checkpoint - CLI commands complete
  - Ensure all tests pass, ask the user if questions arise.

- [x] 10. Cross-compilation and build
  - [x] 10.1 Create Makefile with cross-compilation targets
    - Add `build-all` target compiling for: linux/amd64, linux/arm64, darwin/amd64, darwin/arm64, windows/amd64
    - Add `test` target running `go test ./...`
    - Add `clean` target removing build artifacts
    - Output binaries to `bin/` directory with platform-specific names
    - _Requirements: 6.1_

  - [ ]* 10.2 Write unit tests for cross-compilation verification
    - Verify each platform target compiles without errors
    - Test that binary names follow expected naming convention
    - _Requirements: 6.1_

- [x] 11. Final checkpoint - All tests pass
  - Ensure all tests pass, ask the user if questions arise.

## Notes

- Tasks marked with `*` are optional and can be skipped for faster MVP
- Each task references specific requirements for traceability
- Checkpoints ensure incremental validation
- Property tests validate universal correctness properties from the design document using `pgregory.net/rapid`
- Unit tests validate specific examples and edge cases
- The design specifies Go as the implementation language with `spf13/cobra` for CLI, `charmbracelet/glamour` for rendering, and `pgregory.net/rapid` for property-based tests
- All error messages follow the format: `recall: <description>`
- Files are stored without `.md` extension in the recall directory

## Task Dependency Graph

```json
{
  "waves": [
    { "id": 0, "tasks": ["1.1"] },
    { "id": 1, "tasks": ["2.1", "3.1", "6.1"] },
    { "id": 2, "tasks": ["2.2", "2.3", "2.4", "3.2", "3.3", "4.1"] },
    { "id": 3, "tasks": ["4.2", "5.1"] },
    { "id": 4, "tasks": ["5.2", "5.3"] },
    { "id": 5, "tasks": ["8.1", "8.2", "8.3", "8.4"] },
    { "id": 6, "tasks": ["8.5", "8.6", "8.7", "8.8"] },
    { "id": 7, "tasks": ["10.1"] },
    { "id": 8, "tasks": ["10.2"] }
  ]
}
```
