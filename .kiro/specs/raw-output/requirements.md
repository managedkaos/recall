# Requirements Document

## Introduction

This feature adds a `--raw` (`-r`) flag to the Recall CLI's recall (display) command. When this flag is provided, the application outputs the recall file's content as plain markdown text directly to stdout, bypassing the Glamour markdown renderer. Front-matter stripping is still performed, matching the existing behavior. This flag applies exclusively to the recall (display) operation and does not affect subcommands such as edit, list, or search.

## Glossary

- **Recall_CLI**: The command-line application that provides all recall functionality
- **Recall_File**: A markdown-formatted file stored in the Recall_Directory without a .md file extension
- **Recall_Directory**: The directory on the user's filesystem where recall files are stored
- **Renderer**: The component responsible for parsing markdown content and rendering it with terminal-aware formatting using Glamour
- **Raw_Flag**: The `--raw` / `-r` command-line flag that disables markdown rendering for the recall (display) operation

## Requirements

### Requirement 1: Raw Output Flag Definition

**User Story:** As a user, I want a `--raw` flag on the recall command, so that I can get unformatted markdown output suitable for piping to other tools or reading without ANSI formatting.

#### Acceptance Criteria

1. THE Recall_CLI SHALL accept a `--raw` boolean flag (long form) on the root recall command, defaulting to false when not provided.
2. THE Recall_CLI SHALL accept a `-r` flag (short form) as an alias for `--raw` on the root recall command.
3. WHEN the `--help` flag is provided on the root recall command, THE Recall_CLI SHALL include the `--raw` / `-r` flag in the displayed usage information with a description indicating that it outputs unformatted markdown without ANSI styling.
4. IF both `--raw` and `--edit` flags are provided simultaneously, THEN THE Recall_CLI SHALL reject the command and display an error message indicating that the two flags are mutually exclusive.
5. IF the `--raw` flag is provided without a filename argument, THEN THE Recall_CLI SHALL display the root command help text and exit with a zero exit code.

### Requirement 2: Raw Output Behavior

**User Story:** As a user, I want the raw flag to output plain markdown text without terminal formatting, so that I can pipe recall content to other commands or view it in contexts where ANSI codes are unwanted.

#### Acceptance Criteria

1. WHEN the `--raw` flag is provided with a valid filename argument, THE Recall_CLI SHALL read the corresponding Recall_File, strip the front-matter tags line, and write the remaining content directly to stdout without passing it through the Renderer.
2. WHEN the `--raw` flag is provided with a valid filename argument, THE Recall_CLI SHALL produce output that contains no ANSI escape sequences.
3. WHEN the `--raw` flag is provided with a valid filename argument and the Recall_File contains no front-matter tags line, THE Recall_CLI SHALL write the entire file content to stdout without modification.
4. WHEN the `--raw` flag is provided with a valid filename argument, THE Recall_CLI SHALL exit with a zero exit code after writing the output.

### Requirement 3: Raw Output File-Not-Found Behavior

**User Story:** As a user, I want consistent error behavior regardless of whether raw mode is active, so that my scripts can rely on the same exit code conventions.

#### Acceptance Criteria

1. IF the `--raw` flag is provided with a filename that does not correspond to an existing Recall_File, THEN THE Recall_CLI SHALL exit with exit code 1 without printing any output to stdout or stderr.
2. IF the `--raw` flag is provided and the Recall_Directory cannot be read or does not exist, THEN THE Recall_CLI SHALL exit with exit code 1 without printing any output to stdout or stderr.

### Requirement 4: Raw Flag Scope

**User Story:** As a user, I want the raw flag to only affect the recall display operation, so that subcommands continue to behave as expected without interference.

#### Acceptance Criteria

1. THE Recall_CLI SHALL NOT register the `--raw` or `-r` flag on the "edit" subcommand.
2. THE Recall_CLI SHALL NOT register the `--raw` or `-r` flag on the "list" subcommand.
3. THE Recall_CLI SHALL NOT register the `--raw` or `-r` flag on the "search" subcommand.
4. THE Recall_CLI SHALL NOT register the `--raw` or `-r` flag on the "init" subcommand.
5. THE Recall_CLI SHALL register the `--raw` / `-r` flag as a local (non-persistent) flag on the root recall command so that it is not inherited by subcommands.
6. IF the user provides the `--raw` or `-r` flag on any subcommand, THEN THE Recall_CLI SHALL reject the command with a non-zero exit code and print an error message indicating an unknown flag.
