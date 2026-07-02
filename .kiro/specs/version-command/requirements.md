# Requirements Document

## Introduction

This feature adds a `version` subcommand to the Recall CLI that prints the application's semantic version and exits. The version values (major, minor, patch) are stored in a `version.yml` file in the project root and injected into the Go binary at build time via linker flags. Additionally, "version" is added to the reserved names list to prevent users from creating a recall file named "version".

## Glossary

- **Recall_CLI**: The command-line application (`recall`) that stores, retrieves, and searches markdown reference files.
- **Version_Subcommand**: The `recall version` subcommand that prints version information and exits.
- **Version_File**: A YAML file (`version.yml`) in the project root containing major, minor, and patch version fields.
- **Linker_Flags**: Go build flags (`-ldflags -X`) used to set package-level variable values at compile time.
- **Reserved_Names**: A set of subcommand names that cannot be used as recall filenames.
- **Semantic_Version**: A version string following the major.minor.patch format (e.g., `0.1.0`).

## Requirements

### Requirement 1: Version Subcommand Output

**User Story:** As a user, I want to run `recall version` so that I can see which version of the application is installed.

#### Acceptance Criteria

1. WHEN the user runs `recall version`, THE Recall_CLI SHALL print the version string in the format `recall version <major>.<minor>.<patch>` to standard output.
2. WHEN the user runs `recall version`, THE Recall_CLI SHALL exit with exit code 0.
3. WHEN the user runs `recall version`, THE Recall_CLI SHALL print exactly one line of output followed by a newline.

### Requirement 2: Version File Structure

**User Story:** As a developer, I want the version to be defined in a `version.yml` file so that version values are maintained in a single, human-readable source of truth.

#### Acceptance Criteria

1. THE Version_File SHALL contain three fields: `major`, `minor`, and `patch`.
2. THE Version_File SHALL store each field as an integer value.
3. THE Version_File SHALL reside in the project root directory.

### Requirement 3: Build-Time Version Injection

**User Story:** As a developer, I want the Makefile to inject version values into the binary at build time so that the binary always reflects the current version without manual code edits.

#### Acceptance Criteria

1. WHEN the `build` target is invoked, THE Makefile SHALL read the major, minor, and patch values from the Version_File.
2. WHEN the `build` target is invoked, THE Makefile SHALL pass the version values to the Go compiler using Linker_Flags (`-ldflags -X`).
3. WHEN the `build-all` target is invoked, THE Makefile SHALL pass the version values to the Go compiler using Linker_Flags for each platform build.

### Requirement 4: Package-Level Version Variables

**User Story:** As a developer, I want package-level variables for the version components so that the linker can set them at compile time and the version subcommand can read them.

#### Acceptance Criteria

1. THE Recall_CLI SHALL declare package-level string variables for Version, Major, Minor, and Patch.
2. WHEN Linker_Flags are not provided during build, THE Recall_CLI SHALL default the version variables to empty strings.
3. WHEN the version variables are populated via Linker_Flags, THE Recall_CLI SHALL use those values to construct the Semantic_Version output.

### Requirement 5: Reserved Name Protection

**User Story:** As a user, I want the name "version" to be reserved so that I cannot accidentally create a recall file that conflicts with the version subcommand.

#### Acceptance Criteria

1. THE Recall_CLI SHALL include "version" in the Reserved_Names set.
2. WHEN a user attempts to create or edit a file named "version", THE Recall_CLI SHALL reject the operation with an error message indicating that "version" is a reserved name.

### Requirement 6: Version Subcommand Isolation

**User Story:** As a user, I want the version subcommand to operate independently so that it works without requiring a recall directory or any configuration.

#### Acceptance Criteria

1. WHEN `recall version` is invoked, THE Recall_CLI SHALL print the version without reading or requiring the recall directory.
2. WHEN `recall version` is invoked, THE Recall_CLI SHALL ignore any flags defined on the root command (such as `--raw` or `--edit`).
