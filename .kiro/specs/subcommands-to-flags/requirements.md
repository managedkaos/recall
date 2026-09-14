# Requirements Document

## Introduction

This feature removes the Recall CLI's Cobra subcommands and reimplements each
as a root-level command-line flag. The motivation is that Recall's primary
behavior is to treat the first positional argument as a Recall_File to look up
and render. Subcommands and the positional-filename model are mutually
exclusive namespaces: a file cannot be named `list`, `search`, `edit`, etc.
without colliding with a subcommand. Converting subcommands to flags frees the
filename namespace and makes the "first argument is always a filename" contract
unambiguous.

Where possible, each converted action uses the first letter of the former
subcommand as its single-letter (short) switch and the fully spelled-out
subcommand name as its long switch.

## Glossary

- **Recall_CLI**: The command-line application that provides all recall functionality.
- **Recall_File**: A markdown-formatted file stored in the Recall_Directory without a `.md` extension.
- **Recall_Directory**: The directory where recall files are stored (default `~/.recall` or `RECALL_DIR`).
- **Action_Flag**: A root-level flag that selects a mode of operation other than the default look-up-and-render behavior (`--edit`, `--search`, `--list`, `--init`, `--version`, `--completion`).
- **Default_Recall_Behavior**: When no Action_Flag is provided and a positional argument is present, the Recall_CLI reads the named Recall_File, strips front-matter, and renders it (subject to `--raw`).
- **Edit_Flag**: `--edit` / `-e`.
- **Search_Flag**: `--search` / `-s`.
- **List_Flag**: `--list` / `-l`.
- **Init_Flag**: `--init` / `-i`, with an optional `--init-path <dir>` modifier.
- **Version_Flag**: `--version` / `-v`.
- **Completion_Flag**: `--completion` / `-c`, which takes a shell name value.
- **Raw_Flag**: `--raw` / `-r` (existing; modifies the Default_Recall_Behavior).
- **Tag_Flag**: `--tag` (existing on `list`; becomes a root-level modifier used with List_Flag).

## Requirements

### Requirement 1: Remove Subcommands

**User Story:** As a user, I want no reserved subcommand names, so that any Recall_File name is valid.

#### Acceptance Criteria

1. THE Recall_CLI SHALL NOT register a Cobra subcommand named `edit`.
2. THE Recall_CLI SHALL NOT register a Cobra subcommand named `list` (nor its `ls` alias).
3. THE Recall_CLI SHALL NOT register a Cobra subcommand named `search`.
4. THE Recall_CLI SHALL NOT register a Cobra subcommand named `init`.
5. THE Recall_CLI SHALL NOT register a Cobra subcommand named `version`.
6. THE Recall_CLI SHALL disable or remove Cobra's default `completion` subcommand.
7. WHEN a user runs `recall <name>` where `<name>` was formerly a subcommand (e.g. `list`, `search`), THE Recall_CLI SHALL treat `<name>` as a Recall_File to look up under Default_Recall_Behavior.
8. THE Recall_CLI SHALL remove the reserved-name enforcement (`IsReservedName`) so that no filename is rejected for colliding with a command name.

### Requirement 2: Edit Flag

**User Story:** As a user, I want `-e`/`--edit <name>`, so that I can edit a file without a subcommand.

#### Acceptance Criteria

1. THE Recall_CLI SHALL accept `--edit` (long) and `-e` (short) as a boolean Action_Flag on the root command.
2. WHEN Edit_Flag is provided with a positional argument, THE Recall_CLI SHALL open that Recall_File in `$EDITOR`, creating it if absent, matching the prior `edit` subcommand behavior.
3. IF Edit_Flag is provided without a positional argument, THEN THE Recall_CLI SHALL print an error to stderr and exit with a non-zero code.
4. IF `$EDITOR` is unset when Edit_Flag is used, THEN THE Recall_CLI SHALL print an error to stderr and exit with a non-zero code, matching prior behavior.

### Requirement 3: Search Flag

**User Story:** As a user, I want `-s`/`--search <query>`, so that I can search without a subcommand.

#### Acceptance Criteria

1. THE Recall_CLI SHALL accept `--search` (long) and `-s` (short) as a boolean Action_Flag on the root command.
2. WHEN Search_Flag is provided with a positional query argument, THE Recall_CLI SHALL perform the same case-insensitive substring search and produce the same output format as the prior `search` subcommand.
3. IF Search_Flag is provided without a positional argument, THEN THE Recall_CLI SHALL display help and exit with a zero exit code.

### Requirement 4: List Flag

**User Story:** As a user, I want `-l`/`--list`, so that I can list files without a subcommand.

#### Acceptance Criteria

1. THE Recall_CLI SHALL accept `--list` (long) and `-l` (short) as a boolean Action_Flag on the root command.
2. WHEN List_Flag is provided, THE Recall_CLI SHALL list all Recall_Files sorted alphabetically, matching the prior `list` subcommand output.
3. THE Recall_CLI SHALL accept a `--tag <tag>` modifier that, when combined with List_Flag, filters the listing to files carrying the given tag (case-insensitive), matching prior `list --tag` behavior.
4. WHEN List_Flag is provided, THE Recall_CLI SHALL ignore any positional arguments.

### Requirement 5: Init Flag

**User Story:** As a user, I want `-i`/`--init`, so that I can initialize the Recall_Directory without a subcommand.

#### Acceptance Criteria

1. THE Recall_CLI SHALL accept `--init` (long) and `-i` (short) as a boolean Action_Flag on the root command.
2. WHEN Init_Flag is provided without `--init-path`, THE Recall_CLI SHALL create the resolved Recall_Directory (default or `RECALL_DIR`) if it does not exist and print a confirmation to stdout, matching the prior `init` behavior for the default path.
3. THE Recall_CLI SHALL accept an `--init-path <dir>` modifier that, when combined with Init_Flag, initializes the given directory instead of the default/`RECALL_DIR` location. This replaces the prior optional positional `[path]` argument to `init`.
4. IF `--init-path` is provided without Init_Flag, THEN THE Recall_CLI SHALL print an error to stderr and exit with a non-zero code.

### Requirement 6: Version Flag

**User Story:** As a user, I want `-v`/`--version`, so that I can print the version without a subcommand.

#### Acceptance Criteria

1. THE Recall_CLI SHALL accept `--version` (long) and `-v` (short) as a boolean Action_Flag on the root command.
2. WHEN Version_Flag is provided, THE Recall_CLI SHALL print the build/version metadata and exit with a zero exit code, matching the prior `version` subcommand output.

### Requirement 7: Completion Flag

**User Story:** As a user, I want `-c`/`--completion <shell>`, so that I can generate a shell completion script without a subcommand.

#### Acceptance Criteria

1. THE Recall_CLI SHALL accept `--completion` (long) and `-c` (short) as a string-valued Action_Flag whose value names a shell.
2. THE Recall_CLI SHALL support the shell values `bash`, `zsh`, `fish`, and `powershell`.
3. WHEN Completion_Flag is provided with a supported shell value, THE Recall_CLI SHALL write the corresponding completion script to stdout and exit with a zero exit code.
4. IF Completion_Flag is provided with an unsupported or empty value, THEN THE Recall_CLI SHALL print an error to stderr listing the supported shells and exit with a non-zero code.

### Requirement 8: Single-Letter and Long-Name Switch Consistency

**User Story:** As a user, I want predictable switch names, so that I can remember them.

#### Acceptance Criteria

1. THE Recall_CLI SHALL map each converted action to the first letter of its former subcommand as the short switch: `edit`→`-e`, `search`→`-s`, `list`→`-l`, `init`→`-i`, `version`→`-v`, `completion`→`-c`.
2. THE Recall_CLI SHALL map each converted action to its fully spelled-out former subcommand name as the long switch: `--edit`, `--search`, `--list`, `--init`, `--version`, `--completion`.
3. THE Recall_CLI SHALL preserve Cobra's built-in `-h` / `--help` for help output and SHALL NOT reassign `-h`.

### Requirement 9: Mutual Exclusivity of Action Flags

**User Story:** As a user, I want conflicting actions rejected clearly, so that behavior is never ambiguous.

#### Acceptance Criteria

1. IF more than one Action_Flag is provided in a single invocation, THEN THE Recall_CLI SHALL reject the command with an error to stderr identifying the conflict and exit with exit code 1.
2. THE Recall_CLI SHALL allow Raw_Flag to combine with the Default_Recall_Behavior and SHALL treat Raw_Flag as having no effect when combined with a non-render Action_Flag (or reject it if combined with an incompatible action — to be fixed in design).
3. WHEN no Action_Flag is provided and a positional argument is present, THE Recall_CLI SHALL perform Default_Recall_Behavior.
4. WHEN no Action_Flag and no positional argument are provided, THE Recall_CLI SHALL display help and exit with a zero exit code.

### Requirement 10: Help and Documentation

**User Story:** As a user, I want `--help` and the README to reflect flags, so that I can discover the new interface.

#### Acceptance Criteria

1. WHEN `--help` is requested, THE Recall_CLI SHALL display all Action_Flags with their short and long forms and descriptions, and SHALL NOT list any subcommands.
2. THE README SHALL be updated to document the flag-based interface and remove subcommand usage examples.
