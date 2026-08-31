# Requirements Document

## Introduction

This feature adds a `--search` (`-s`) flag to the Recall CLI's root command. When this flag is provided with a query argument, the application performs the same case-insensitive substring search across all recall files as the existing `search` subcommand, and prints identical output. This gives users a shorthand — `recall -s "query"` — equivalent to `recall search "query"`, mirroring how the `--edit` (`-e`) and `--raw` (`-r`) flags are wired on the root command. The flag applies exclusively to the root recall command and does not affect subcommands.

## Glossary

- **Recall_CLI**: The command-line application that provides all recall functionality
- **Recall_File**: A markdown-formatted file stored in the Recall_Directory without a `.md` file extension
- **Recall_Directory**: The directory on the user's filesystem where recall files are stored (default `~/.recall` or `RECALL_DIR`)
- **Search_Engine**: The `internal/search` package that performs a case-insensitive substring scan across all Recall_Files and returns matches grouped by file
- **Search_Subcommand**: The existing `recall search <query>` subcommand
- **Search_Flag**: The `--search` / `-s` command-line flag that triggers search behavior on the root recall command
- **Edit_Flag**: The existing `--edit` / `-e` flag on the root recall command
- **Raw_Flag**: The existing `--raw` / `-r` flag on the root recall command

## Requirements

### Requirement 1: Search Flag Definition

**User Story:** As a user, I want a `--search` flag on the recall command, so that I can search all recall files without typing the `search` subcommand.

#### Acceptance Criteria

1. THE Recall_CLI SHALL accept a `--search` boolean flag (long form) on the root recall command, defaulting to false when not provided.
2. THE Recall_CLI SHALL accept a `-s` flag (short form) as an alias for `--search` on the root recall command.
3. WHEN the `--help` flag is provided on the root recall command, THE Recall_CLI SHALL include the `--search` / `-s` flag in the displayed usage information with a description indicating that it searches all recall files for the given query.
4. THE Recall_CLI SHALL register the `--search` / `-s` flag as a local (non-persistent) flag on the root recall command so that it is not inherited by subcommands.

### Requirement 2: Search Flag Behavior

**User Story:** As a user, I want `recall -s "query"` to behave exactly like `recall search "query"`, so that the flag is a reliable shorthand.

#### Acceptance Criteria

1. WHEN the `--search` flag is provided with a query argument, THE Recall_CLI SHALL treat the first positional argument as the search query and invoke the Search_Engine against all Recall_Files in the Recall_Directory.
2. WHEN the `--search` flag is provided with a query argument, THE Recall_CLI SHALL produce output identical to the output of the Search_Subcommand invoked with the same query and the same Recall_Directory contents.
3. WHEN the `--search` flag is provided with a query argument that matches one or more lines, THE Recall_CLI SHALL print each match in the form `filename:linenumber:linecontent`, grouping matches by file and separating file groups with a line containing `----------`.
4. WHEN the `--search` flag is provided with a query argument that matches nothing, THE Recall_CLI SHALL produce no output and exit with a zero exit code.
5. WHEN the `--search` flag is provided with a query argument, THE Recall_CLI SHALL exit with a zero exit code after a successful search.

### Requirement 3: Search Flag Argument Handling

**User Story:** As a user, I want predictable behavior when I forget the query, so that the tool guides me rather than failing obscurely.

#### Acceptance Criteria

1. IF the `--search` flag is provided without any positional argument, THEN THE Recall_CLI SHALL display the root command help text and exit with a zero exit code.

### Requirement 4: Mutual Exclusivity

**User Story:** As a user, I want conflicting flags to be rejected clearly, so that I do not get ambiguous behavior.

#### Acceptance Criteria

1. IF both `--search` and `--edit` flags are provided simultaneously, THEN THE Recall_CLI SHALL reject the command and display an error message to stderr indicating that the two flags are mutually exclusive, and exit with exit code 1.
2. THE Recall_CLI SHALL allow the `--search` and `--raw` flags to be provided together without error, since the Raw_Flag has no effect on search output.

### Requirement 5: Search Flag Scope

**User Story:** As a user, I want the search flag to only affect the root recall command, so that subcommands continue to behave as expected.

#### Acceptance Criteria

1. THE Recall_CLI SHALL NOT register the `--search` or `-s` flag on the "edit" subcommand.
2. THE Recall_CLI SHALL NOT register the `--search` or `-s` flag on the "list" subcommand.
3. THE Recall_CLI SHALL NOT register the `--search` or `-s` flag on the "search" subcommand.
4. THE Recall_CLI SHALL NOT register the `--search` or `-s` flag on the "init" subcommand.
5. THE Recall_CLI SHALL NOT register the `--search` or `-s` flag on the "version" subcommand.
6. IF the user provides the `--search` or `-s` flag on any subcommand, THEN THE Recall_CLI SHALL reject the command with a non-zero exit code and print an error message indicating an unknown flag.

### Requirement 6: Error and Exit-Code Parity

**User Story:** As a user, I want the search flag to fail the same way the search subcommand does, so that my scripts behave consistently.

#### Acceptance Criteria

1. IF the Recall_Directory cannot be resolved or created while handling the `--search` flag, THEN THE Recall_CLI SHALL print an error to stderr and exit with exit code 1, matching the Search_Subcommand's behavior.
2. IF the Search_Engine returns an error while handling the `--search` flag, THEN THE Recall_CLI SHALL print an error to stderr and exit with exit code 1, matching the Search_Subcommand's behavior.
