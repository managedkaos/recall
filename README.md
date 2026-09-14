# Recall

A command-line tool for storing and recalling markdown-formatted reference files.

## Build Requirements

- Go 1.22 or later
- Make (optional, for cross-compilation targets)

## Building

Build for your current platform:

```bash
make build
```

The binary is output to `bin/recall`.

To cross-compile for all supported platforms (linux/amd64, linux/arm64, darwin/amd64, darwin/arm64, windows/amd64):

```bash
make build-all
```

Alternatively, build directly with Go:

```bash
go build -o recall .
```

## Running

Recall's primary action is to look up a file: the first argument is always
treated as the name of a file to display. All other operations are selected
with root-level flags. Each action has a single-letter short form and a
spelled-out long form. Only one action flag may be used at a time.

| Action | Short | Long |
| --- | --- | --- |
| Edit a file in `$EDITOR` | `-e` | `--edit` |
| Search all files | `-s` | `--search` |
| List all files | `-l` | `--list` |
| Initialize the directory | `-i` | `--init` |
| Print the version | `-v` | `--version` |
| Generate a completion script | `-c <shell>` | `--completion <shell>` |

### Recall a file

Display a stored file with terminal-rendered markdown:

```bash
recall docker
```

Output the unformatted markdown with `-r` / `--raw`:

```bash
recall -r docker
```

Because there are no subcommands, any name is a valid filename — including
names like `list` or `search`.

### Store a file

Create or edit a recall file using your `$EDITOR`:

```bash
recall -e docker
```

This opens a file named `docker` in your editor. Write your notes in markdown
format and save. If the file doesn't exist, it is created first.

### List files

List all stored recall files:

```bash
recall -l
```

Filter by tag:

```bash
recall -l --tag devops
```

### Search

Search across all files for a string (case-insensitive):

```bash
recall -s "docker compose"
```

### Initialize (optional)

Create the recall directory explicitly. This is optional — the directory is
created automatically on first use.

```bash
recall -i
```

By default this initializes `~/.recall` (or the path from `RECALL_DIR`). To
initialize a specific directory, use `--init-path`:

```bash
recall -i --init-path ~/notes/recall
```

### Version

```bash
recall -v
```

### Shell completion

Generate a completion script for your shell (`bash`, `zsh`, `fish`, or
`powershell`):

```bash
recall -c bash > /etc/bash_completion.d/recall
```

### Help

```bash
recall --help
```

## Configuration

Recall stores files in `~/.recall` by default. Override this by setting the `RECALL_DIR` environment variable:

```bash
export RECALL_DIR=~/notes/recall
```

The editor is determined by the `$EDITOR` environment variable:

```bash
export EDITOR=vim
```

## File Format

Recall files are plain text with markdown content and no `.md` extension. Optionally, the first line can contain tags:

```
tags: docker, devops, containers
# Docker Cheatsheet

## Running Containers
- `docker run -d --name myapp nginx`
- `docker ps -a`
```

Tags enable filtering with `recall -l --tag <tag>`.

## Supported Platforms

- Linux (amd64, arm64)
- macOS (amd64, arm64)
- Windows (amd64)
