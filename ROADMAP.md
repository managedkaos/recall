# Roadmap

Planned features and improvements for `recall`. Items are aspirational and
subject to change.

## Features to add

### Search via `-s` flag

Add a root-level `-s` flag as an alternative to the existing `search`
subcommand, mirroring how the `-e` (edit) and `-r` (raw) flags are wired on
the root command. This lets you run `recall -s "docker compose"` alongside
the current `recall search "docker compose"`.

### Encrypt and decrypt files

Support encrypting and decrypting individual recall files at rest, built
directly into Recall using the [`age`](https://pkg.go.dev/filippo.io/age)
library — no external tools required. Two models are planned: application-managed
keys (default) and passphrase-based encryption (portable alternative).

See the full design in
[`.kiro/specs/file-encryption/`](.kiro/specs/file-encryption/design.md).

### Metrics collection

Track usage metrics for recall operations — files edited/created, recalled,
encrypted, and decrypted, among others — to give insight into how the tool
is used over time.

### TUI-based manager

Provide an interactive terminal UI to browse stored files and view collected
metrics, with the possibility of editing files directly from the interface.
