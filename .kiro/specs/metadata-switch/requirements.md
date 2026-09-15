# Requirements Document

## Introduction

This feature adds a `--metadata` (`-m`) flag to the Recall CLI's root command. When this flag is provided with one or more filename arguments, the application reports, for each named Recall_File that exists, its file path, filesystem metadata (size, modification time, and — where the platform exposes them — creation and change times), and any tags parsed from the file's front-matter. An optional sub-flag `--json` (`-j`) switches the output from human-readable text to a machine-readable JSON document. The flag applies exclusively to the root recall command and is not inherited by subcommands, mirroring how `--edit`, `--search`, `--list`, and `--raw` are wired.

## Glossary

- **Recall_CLI**: The command-line application that provides all recall functionality.
- **Recall_File**: A markdown-formatted file stored in the Recall_Directory without a `.md` file extension.
- **Recall_Directory**: The directory on the user's filesystem where recall files are stored (default `~/.recall` or `RECALL_DIR`).
- **Metadata_Flag**: The `--metadata` / `-m` command-line flag that triggers metadata-reporting behavior on the root recall command.
- **JSON_Flag**: The `--json` / `-j` sub-flag that changes Metadata_Flag output to JSON format.
- **Frontmatter_Parser**: The `internal/frontmatter` package that extracts tags from the first line of a Recall_File's content.
- **File_Metadata**: The set of reported attributes for a Recall_File — absolute path, byte size, modification time, and, when available on the platform, creation time and change (inode-change) time.
- **Edit_Flag**: The existing `--edit` / `-e` flag on the root recall command.
- **Search_Flag**: The existing `--search` / `-s` flag on the root recall command.
- **List_Flag**: The existing `--list` / `-l` flag on the root recall command.
- **Raw_Flag**: The existing `--raw` / `-r` flag on the root recall command.

## Requirements

### Requirement 1: Metadata Flag Definition

**User Story:** As a user, I want a `--metadata` flag on the recall command, so that I can inspect the file details and tags of stored recall files.

#### Acceptance Criteria

1. THE Recall_CLI SHALL accept a `--metadata` boolean flag (long form) on the root recall command, defaulting to false when not provided.
2. THE Recall_CLI SHALL accept a `-m` flag (short form) as an alias for `--metadata` on the root recall command.
3. THE Recall_CLI SHALL accept a `--json` boolean flag (long form) as a modifier of `--metadata`, defaulting to false when not provided.
4. THE Recall_CLI SHALL accept a `-j` flag (short form) as an alias for `--json`.
5. WHEN the `--help` flag is provided on the root recall command, THE Recall_CLI SHALL include the `--metadata` / `-m` and `--json` / `-j` flags in the displayed usage information with descriptions.
6. THE Recall_CLI SHALL register the `--metadata` / `-m` and `--json` / `-j` flags as local (non-persistent) flags on the root recall command so that they are not inherited by subcommands.

### Requirement 2: Metadata Reporting Behavior

**User Story:** As a user, I want `recall -m <file>` to report the path, filesystem metadata, and tags for each named file, so that I can understand a file's provenance at a glance.

#### Acceptance Criteria

1. WHEN the `--metadata` flag is provided with one or more filename arguments, THE Recall_CLI SHALL process each argument in the order given on the command line.
2. WHEN a named argument corresponds to a Recall_File that exists in the Recall_Directory, THE Recall_CLI SHALL report the file's absolute path.
3. WHEN a named argument corresponds to a Recall_File that exists, THE Recall_CLI SHALL report the file's size in bytes and its last modification time.
4. WHEN a named argument corresponds to a Recall_File that exists AND the underlying platform exposes creation time and/or inode-change time, THE Recall_CLI SHALL report those times; WHERE the platform does not expose a given time, THE Recall_CLI SHALL omit it rather than reporting a placeholder or erroring.
5. WHEN a named argument corresponds to a Recall_File that exists, THE Recall_CLI SHALL parse the file's front-matter with the Frontmatter_Parser and report any tags found; WHERE the file has no tags, THE Recall_CLI SHALL report an empty tag set.
6. WHEN the `--metadata` flag is provided with one filename argument for an existing file, THE Recall_CLI SHALL exit with a zero exit code after reporting.

### Requirement 3: Missing and Invalid File Handling

**User Story:** As a user, I want clear behavior when a named file does not exist, so that I can distinguish present from absent files.

#### Acceptance Criteria

1. IF a named argument does not correspond to an existing Recall_File, THEN THE Recall_CLI SHALL print an error to stderr identifying the missing file and SHALL NOT print File_Metadata for that argument.
2. WHEN the `--metadata` flag is provided with multiple filename arguments where some exist and some do not, THE Recall_CLI SHALL report File_Metadata for every existing file and emit a stderr error for every missing file.
3. IF every named argument is missing, THEN THE Recall_CLI SHALL exit with a non-zero exit code.
4. IF at least one named argument corresponds to an existing file, THEN THE Recall_CLI SHALL exit with a zero exit code regardless of other missing arguments.
5. IF the `--metadata` flag is provided without any positional argument, THEN THE Recall_CLI SHALL display the root command help text and exit with a zero exit code.

### Requirement 4: JSON Output

**User Story:** As a user, I want `recall -m -j <file>` to emit JSON, so that I can consume metadata in scripts and other tools.

#### Acceptance Criteria

1. WHEN both the `--metadata` and `--json` flags are provided, THE Recall_CLI SHALL emit the File_Metadata for all existing named files as a single well-formed JSON document to stdout.
2. WHEN emitting JSON, THE Recall_CLI SHALL represent the result as a JSON array of objects, one object per existing named file, each containing at least the file name, absolute path, size in bytes, modification time, and tags.
3. WHERE a platform-specific time (creation or change time) is unavailable, THE Recall_CLI SHALL omit that field from the JSON object rather than emitting a null or placeholder value.
4. WHEN emitting JSON, THE Recall_CLI SHALL format timestamps as RFC 3339 strings.
5. WHEN emitting JSON and no named file exists, THE Recall_CLI SHALL emit an empty JSON array `[]` to stdout and report missing files to stderr.
6. IF the `--json` flag is provided without the `--metadata` flag, THEN THE Recall_CLI SHALL reject the command with an error to stderr indicating that `--json` requires `--metadata`, and exit with exit code 1.

### Requirement 5: Mutual Exclusivity

**User Story:** As a user, I want conflicting action flags to be rejected clearly, so that I do not get ambiguous behavior.

#### Acceptance Criteria

1. THE Recall_CLI SHALL treat `--metadata` as an action flag subject to the existing "only one action flag at a time" rule alongside `--edit`, `--search`, `--list`, `--init`, `--version`, and `--completion`.
2. IF the `--metadata` flag is provided together with any other action flag, THEN THE Recall_CLI SHALL reject the command with an error to stderr and exit with exit code 1.
3. THE Recall_CLI SHALL allow the `--metadata` and `--json` flags to be provided together, as `--json` is a modifier of `--metadata`.

### Requirement 6: Metadata Flag Scope

**User Story:** As a user, I want the metadata flags to only affect the root recall command, so that subcommands continue to behave as expected.

#### Acceptance Criteria

1. THE Recall_CLI SHALL NOT register the `--metadata` / `-m` or `--json` / `-j` flags on any subcommand.
2. IF the user provides the `--metadata`, `-m`, `--json`, or `-j` flag on any subcommand, THEN THE Recall_CLI SHALL reject the command with a non-zero exit code and print an error message indicating an unknown flag.

### Requirement 7: Read-Only Guarantee

**User Story:** As a user, I want metadata reporting to never change my files, so that inspection is safe.

#### Acceptance Criteria

1. WHEN the `--metadata` flag is provided (with or without `--json`), THE Recall_CLI SHALL NOT create, modify, or delete any Recall_File.
2. IF the Recall_Directory cannot be resolved or created while handling the `--metadata` flag, THEN THE Recall_CLI SHALL print an error to stderr and exit with exit code 1.
