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

### Initialize (optional)

Create the recall directory explicitly. This is optional — the directory is created automatically on first use.

```bash
recall init
```

### Store a file

Create or edit a recall file using your `$EDITOR`:

```bash
recall edit docker
```

This opens a file named `docker` in your editor. Write your notes in markdown format and save.

### Recall a file

Display a stored file with terminal-rendered markdown:

```bash
recall docker
```

You can also use the `-e` flag to quickly edit:

```bash
recall -e docker
```

### List files

List all stored recall files:

```bash
recall list
```

Filter by tag:

```bash
recall list --tag devops
```

### Search

Search across all files for a string (case-insensitive):

```bash
recall search "docker compose"
```

### Help

```bash
recall --help
recall edit --help
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

Tags enable filtering with `recall list --tag <tag>`.

## Supported Platforms

- Linux (amd64, arm64)
- macOS (amd64, arm64)
- Windows (amd64)
