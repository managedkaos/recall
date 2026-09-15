# Requirements: Version Metadata

The CLI prints version and build metadata through `--version` / `-v` and exits
successfully without requiring a recall directory.

- Use one build-time string, `cmd.Version`, injected through Go linker flags.
- Local `make build` and `make build-all` derive this string from the nearest
  reachable numeric, v-prefixed numeric, or V-prefixed numeric Git tag. Include
  commit distance and hash between tags and `-local` for tracked changes.
- Support annotated, lightweight, prerelease, and date-based tags.
- Allow environment and command-line `VERSION` overrides, with command-line
  precedence. Overrides bypass Git derivation.
- Missing Git, repository metadata, or matching tags must yield `unknown`.
- Trim whitespace and one leading v/V for display; preserve all suffixes.
  Empty injected values yield `unknown`, including plain Go builds.
- GoReleaser supplies release and snapshot versions independently of Make.
- CI Make builds fetch full history and tags. Releases do not synchronize a
  separate version file or open version synchronization pull requests.
- Preserve all existing build metadata fields and both version flags.
