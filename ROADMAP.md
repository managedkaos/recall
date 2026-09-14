# Roadmap

Planned features and improvements for `recall`. Items are aspirational and
subject to change.

## Features to add

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
