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
| Report file metadata | `-m` | `--metadata` |
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

### Report metadata

Report the file path, filesystem metadata (size, modification time, and — where
the platform exposes them — creation time), and any tags for one or
more files:

```bash
recall -m docker
```

You can pass multiple names; metadata is reported for each file that exists, and
a message is printed to stderr for any that do not. The command exits with a
non-zero status only when none of the named files exist.

```bash
recall -m docker kubernetes notes
```

Use `-j` / `--json` to emit the report as a JSON array instead of text. This
modifier is only valid together with `--metadata`:

```bash
recall -m -j docker
```

Timestamps in JSON output are formatted as RFC 3339. Note that the available
timestamps vary by platform: the modification time is reported on all
platforms; macOS and Windows also report the creation time, while Linux reports
only the modification time. Unavailable times are omitted rather than reported
as empty.

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

Once the script is installed, pressing `<TAB>` completes the names of your
stored recall files. Completion is dynamic — it reads the current recall
directory (honoring `RECALL_DIR`) each time, so newly created files are
available immediately without regenerating the script:

```bash
recall doc<TAB>        # completes to stored files beginning with "doc"
recall -e doc<TAB>     # same completion when editing
recall -m doc<TAB>     # metadata accepts multiple names; keeps completing
```

Completion is context-aware:

- The default lookup and `--edit`, `--raw`, and `--metadata` complete stored
  file names.
- `--completion` / `-c` completes the shell name (`bash`, `zsh`, `fish`,
  `powershell`).
- `--search` takes a free-text query, and `--list`, `--init`, and `--version`
  take no filename, so none of these offer file-name candidates.

Install the script for your shell in the usual location. For example:

```bash
# bash (Linux)
recall -c bash | sudo tee /etc/bash_completion.d/recall > /dev/null

# zsh (place on your $fpath, e.g. ~/.zsh/completions/_recall)
recall -c zsh > ~/.zsh/completions/_recall

# fish
recall -c fish > ~/.config/fish/completions/recall.fish
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

## Releasing

Releases are produced by [GoReleaser](https://goreleaser.com). **The git tag is
the authoritative version.** To cut a release, push a semantic-version tag (or
run the "Build and Release Binary" workflow manually with a tag input):

```bash
git tag v1.2.3
git push origin v1.2.3
```

The release workflow then:

1. **Builds** all five platform archives with GoReleaser, injecting the tag as
   the binary version (`recall -v` reports `1.2.3`).
2. **Smoke-tests** every archive on its native runner by extracting it and
   running `recall --version`.
3. **Publishes** a GitHub release with the archives, a `checksums.txt`, and
   auto-generated release notes — only after the smoke tests pass.
4. **Syncs `version.yml`** by opening a pull request that sets `major`, `minor`,
   and `patch` to match the tag. Because `main` is protected, this arrives as a
   PR for review rather than a direct push.

`version.yml` remains the source of truth for **local** `make` builds, which
compose the version from its components. Release builds prefer the tag-injected
version, so the two stay consistent via the sync PR.

### Local release dry run

Build the full set of release archives locally (no publishing) with GoReleaser:

```bash
make snapshot
```

This requires GoReleaser to be installed and writes artifacts to `dist/`.

## Supported Platforms

- Linux (amd64, arm64)
- macOS (amd64, arm64)
- Windows (amd64)
