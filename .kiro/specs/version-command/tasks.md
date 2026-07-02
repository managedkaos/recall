# Implementation Plan: Version Command

## Overview

Add a `recall version` subcommand that prints the semantic version string and exits. The version values are stored in `version.yml` at the project root and injected into the binary at build time via Go linker flags. The implementation involves creating the version file, adding a new cobra command, updating the Makefile, registering "version" as a reserved name, and writing tests.

## Tasks

- [ ] 1. Create version file and version command source
  - [ ] 1.1 Create `version.yml` in the project root
    - Add the file with fields: `major: 0`, `minor: 1`, `patch: 0`
    - _Requirements: 2.1, 2.2, 2.3_

  - [ ] 1.2 Create `cmd/version.go` with the version subcommand
    - Declare package-level variables: `Version`, `Major`, `Minor`, `Patch` (all strings, default empty)
    - Implement `formatVersion(major, minor, patch string) string` that returns `"unknown"` if any component is empty, otherwise returns `major + "." + minor + "." + patch`
    - Define `versionCmd` as a cobra command with `Use: "version"`, `Args: cobra.NoArgs`, and `RunE: runVersion`
    - Implement `runVersion` to print `recall version <formatted>` and return nil
    - Register the command via `init()` calling `rootCmd.AddCommand(versionCmd)`
    - _Requirements: 1.1, 1.2, 1.3, 4.1, 4.2, 4.3, 6.1, 6.2_

- [ ] 2. Update existing files and write tests
  - [ ] 2.1 Add "version" to `reservedNames` in `cmd/root.go`
    - Add `"version": true` to the `reservedNames` map
    - _Requirements: 5.1, 5.2_

  - [ ] 2.2 Update `Makefile` to extract version from `version.yml` and inject via ldflags
    - Add variables: `VERSION_MAJOR`, `VERSION_MINOR`, `VERSION_PATCH` extracted via `grep` and `awk`
    - Add `LDFLAGS` variable with `-X $(MODULE)/cmd.Major=...` for each component
    - Update the `build` target to include `$(LDFLAGS)`
    - Update the `build-all` target to include `$(LDFLAGS)` for each platform build
    - _Requirements: 3.1, 3.2, 3.3_

  - [ ]* 2.3 Write unit tests for `formatVersion` in `cmd/version_test.go`
    - Test that `formatVersion("0", "1", "0")` returns `"0.1.0"`
    - Test that `formatVersion("", "1", "0")` returns `"unknown"`
    - Test that `formatVersion("", "", "")` returns `"unknown"`
    - Test that the version command produces exactly one line of output
    - Test that `IsReservedName("version")` returns `true`
    - _Requirements: 1.1, 1.3, 4.2, 5.1_

  - [ ]* 2.4 Write property test for format correctness
    - **Property 1: Version format correctness**
    - For any three non-empty strings, `formatVersion` returns `major + "." + minor + "." + patch`
    - Use `pgregory.net/rapid` with minimum 100 iterations
    - **Validates: Requirements 1.1, 1.3, 4.3**

  - [ ]* 2.5 Write property test for empty component detection
    - **Property 2: Empty components produce "unknown"**
    - For any combination where at least one component is empty, `formatVersion` returns `"unknown"`
    - Use `pgregory.net/rapid` with minimum 100 iterations
    - **Validates: Requirements 4.2**

- [ ] 3. Checkpoint - Ensure all tests pass
  - Ensure all tests pass, ask the user if questions arise.

- [ ] 4. Final wiring and verification
  - [ ] 4.1 Verify end-to-end: `make build` then `./bin/recall version` outputs expected format
    - Confirm that the built binary outputs `recall version 0.1.0`
    - Confirm exit code is 0
    - Confirm no recall directory or config is required
    - _Requirements: 1.1, 1.2, 3.1, 3.2, 6.1_

- [ ] 5. Final checkpoint - Ensure all tests pass
  - Ensure all tests pass, ask the user if questions arise.

## Notes

- Tasks marked with `*` are optional and can be skipped for faster MVP
- Each task references specific requirements for traceability
- Checkpoints ensure incremental validation
- Property tests validate universal correctness properties from the design document
- Unit tests validate specific examples and edge cases
- The design uses Go directly, so all implementation uses Go with `pgregory.net/rapid` for property tests

## Task Dependency Graph

```json
{
  "waves": [
    { "id": 0, "tasks": ["1.1", "1.2"] },
    { "id": 1, "tasks": ["2.1", "2.2"] },
    { "id": 2, "tasks": ["2.3", "2.4", "2.5"] },
    { "id": 3, "tasks": ["4.1"] }
  ]
}
```
