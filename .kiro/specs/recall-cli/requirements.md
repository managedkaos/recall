# Requirements Document

## Introduction

Recall is a command-line application written in Go that serves as a replacement for the deprecated "cheat" application. Recall allows users to store, retrieve, edit, list, and search small pieces of information (stored as markdown files without the .md extension) from the command line. The application renders markdown content with terminal-aware formatting and supports cross-compilation for Linux, macOS, and Windows.

## Glossary

- **Recall_CLI**: The command-line application that provides all recall functionality
- **Recall_Directory**: The directory on the user's filesystem where recall files are stored; defaults to "$HOME/.config/recall" and is configurable via an environment variable or the init command
- **Recall_File**: A markdown-formatted file stored in the Recall_Directory without a .md file extension, named after the topic it describes
- **Renderer**: The component responsible for parsing markdown content and rendering it with terminal-aware formatting (headings, lists, emphasis, etc.)
- **Editor**: The external text editor identified by the user's $EDITOR environment variable
- **Search_Engine**: The component responsible for finding string matches across all Recall_Files

## Requirements

### Requirement 1: Recall a File

**User Story:** As a user, I want to recall information about a topic by name, so that I can quickly access stored details from the command line.

#### Acceptance Criteria

1. WHEN a filename argument is provided without a subcommand, THE Recall_CLI SHALL perform a case-sensitive exact match against filenames in the Recall_Directory and render the corresponding Recall_File contents to standard output.
2. IF a filename argument is provided and no matching Recall_File exists in the Recall_Directory, THEN THE Recall_CLI SHALL exit with a non-zero exit code without printing any output to standard output or standard error.
3. WHEN a Recall_File is rendered, THE Renderer SHALL apply distinct visual formatting for markdown headings (H1 through H6), ordered lists, unordered lists, bold text, italic text, inline code, and fenced code blocks using ANSI escape sequences supported by the terminal.
4. WHEN a Recall_File contains a front-matter tags line, THE Renderer SHALL exclude the front-matter line and render only the remaining content.

### Requirement 2: Edit a File

**User Story:** As a user, I want to create or edit recall files using my preferred text editor, so that I can easily add and update stored information.

#### Acceptance Criteria

1. WHEN the "edit" subcommand or the "-e" flag is provided with a filename argument, THE Recall_CLI SHALL launch the Editor as a subprocess with the full path to the corresponding Recall_File as an argument and wait for the Editor process to exit before returning control to the terminal.
2. WHEN the "edit" subcommand or the "-e" flag is provided with a filename that does not correspond to an existing Recall_File, THE Recall_CLI SHALL create a new empty file with that name in the Recall_Directory and open it in the Editor.
3. IF the $EDITOR environment variable is not set, THEN THE Recall_CLI SHALL exit with a non-zero exit code and print an error message to standard error indicating that the $EDITOR variable is not configured.
4. IF the Editor process exits with a non-zero exit code, THEN THE Recall_CLI SHALL exit with a non-zero exit code.

### Requirement 3: List Files

**User Story:** As a user, I want to list all available recall files, so that I can see what information I have stored.

#### Acceptance Criteria

1. WHEN the "list" subcommand is provided, THE Recall_CLI SHALL print the names of all Recall_Files in the Recall_Directory to standard output, one filename per line, sorted in ascending alphabetical order, excluding subdirectories and hidden files.
2. WHEN the "list" subcommand is provided and the Recall_Directory contains no Recall_Files, THE Recall_CLI SHALL produce no output and exit with a zero exit code.
3. IF the Recall_Directory cannot be read when the "list" subcommand is provided, THEN THE Recall_CLI SHALL print a descriptive error message to standard error and exit with a non-zero exit code.

### Requirement 4: Search Files

**User Story:** As a user, I want to search across all stored files for a specific string, so that I can find information when I don't remember which file contains it.

#### Acceptance Criteria

1. WHEN the "search" subcommand is provided with a search string, THE Search_Engine SHALL perform a case-insensitive substring scan of all Recall_Files in the Recall_Directory, matching any line that contains the search string regardless of letter case.
2. WHEN a match is found, THE Recall_CLI SHALL print the result in the format: filename, followed by a colon, followed by the line number (starting from 1), followed by a colon, followed by the content of the matching line.
3. WHEN matches are found in a single Recall_File, THE Recall_CLI SHALL print matching lines in ascending line-number order.
4. WHEN matches are found in multiple Recall_Files, THE Recall_CLI SHALL separate each file's group of results with a line of ten dash characters ("----------"), with each file's matches printed in ascending line-number order.
5. WHEN no matches are found across any Recall_Files, THE Recall_CLI SHALL produce no output and exit with a zero exit code.
6. IF the "search" subcommand is provided without a search string argument, THEN THE Recall_CLI SHALL print an error message to standard error and exit with a non-zero exit code.

### Requirement 5: Configure Storage Location

**User Story:** As a user, I want to configure where my recall files are stored, so that I can organize my filesystem according to my preferences.

#### Acceptance Criteria

1. IF the RECALL_DIR environment variable is set to a non-empty value, THEN THE Recall_CLI SHALL use that value as the Recall_Directory path.
2. IF the RECALL_DIR environment variable is not set or is set to an empty string, THEN THE Recall_CLI SHALL default the Recall_Directory to "$HOME/.config/recall".
3. IF the configured Recall_Directory does not exist, THEN THE Recall_CLI SHALL create the directory, including any missing intermediate directories, before performing any file operations.
4. IF the configured Recall_Directory cannot be created or is not writable, THEN THE Recall_CLI SHALL print a descriptive error message to standard error and exit with a non-zero exit code.

### Requirement 6: Cross-Platform Compilation

**User Story:** As a user, I want to use Recall on Linux, macOS, or Windows, so that I can access my stored information regardless of operating system.

#### Acceptance Criteria

1. THE Recall_CLI SHALL compile without error for linux/amd64, linux/arm64, darwin/amd64, darwin/arm64, and windows/amd64 target platforms.
2. THE Recall_CLI SHALL resolve the user's home directory using platform-appropriate methods on each supported operating system, returning an absolute path containing no platform-incompatible separators.
3. IF the user's home directory cannot be determined at runtime, THEN THE Recall_CLI SHALL exit with a non-zero exit code and print a descriptive error message to standard error.
4. THE Recall_CLI SHALL construct all Recall_Directory and Recall_File paths using the operating system's native path separator so that file operations succeed on each supported platform.

### Requirement 7: Help and Usage Information

**User Story:** As a user, I want to see usage information for the application, so that I can learn the available commands and flags.

#### Acceptance Criteria

1. WHEN the "--help" or "-h" flag is provided, THE Recall_CLI SHALL print usage information to standard output that includes the application name, a list of all available subcommands each with a brief description, and a list of all available flags each with a brief description, and SHALL exit with a zero exit code.
2. IF an unrecognized subcommand or flag is provided, THEN THE Recall_CLI SHALL print an error message indicating the unrecognized input followed by the usage information to standard error and exit with a non-zero exit code.
3. WHEN the "--help" or "-h" flag is provided along with a valid subcommand, THE Recall_CLI SHALL print usage information specific to that subcommand, including its accepted arguments and flags, to standard output and exit with a zero exit code.

### Requirement 8: Tag-Based Organization

**User Story:** As a user, I want to tag my recall files with categories, so that I can organize and filter information by topic.

#### Acceptance Criteria

1. WHEN a Recall_File contains a first line in the format "tags: tag1, tag2, tag3", THE Recall_CLI SHALL parse the comma-separated values as tags associated with the file, trimming leading and trailing whitespace from each tag value and ignoring any empty values resulting from consecutive commas or trailing commas.
2. WHEN the "list" subcommand is provided with a "--tag" flag and a tag value, THE Recall_CLI SHALL print only the names of Recall_Files whose parsed tags contain a case-insensitive match of the specified tag value, one name per line.
3. WHEN a Recall_File is rendered, THE Renderer SHALL exclude the front-matter tags line from the displayed output.
4. WHEN the "list" subcommand is provided with a "--tag" flag and no Recall_Files contain a matching tag, THE Recall_CLI SHALL produce no output and exit with a zero exit code.

### Requirement 9: Initialize Recall Directory

**User Story:** As a user, I want to explicitly initialize the recall storage directory, so that I can confirm or customize where my recall files are stored.

#### Acceptance Criteria

1. WHEN the "init" subcommand is provided without a path argument, THE Recall_CLI SHALL create the default Recall_Directory (including any missing intermediate directories) and print a confirmation message to standard output indicating which directory was initialized.
2. WHEN the "init" subcommand is provided with a path argument, THE Recall_CLI SHALL create the directory at the specified path (including any missing intermediate directories), persist that path as the configured Recall_Directory location by writing it to a configuration file, and print a confirmation message to standard output indicating which directory was initialized.
3. WHEN the "init" subcommand is provided and the target directory already exists, THE Recall_CLI SHALL print a message to standard output informing the user the directory already exists and exit with a zero exit code.
4. IF the target directory cannot be created due to insufficient permissions or an invalid path, THEN THE Recall_CLI SHALL print a descriptive error message to standard error and exit with a non-zero exit code.
5. WHEN a custom path has been persisted via the "init" subcommand, THE Recall_CLI SHALL use that persisted path as the Recall_Directory for all subsequent operations unless overridden by the RECALL_DIR environment variable.
