# Implementation Plan: File Encryption

## Overview

Add local file encryption and decryption to the Recall CLI using the `filippo.io/age` library. Support application-managed key encryption (default) and passphrase-based encryption (optional). Introduce an `internal/crypto` package and `cmd/encrypt.go` / `cmd/decrypt.go` (and optionally `cmd/keygen.go`), following existing Cobra patterns. Keep identity generation separate from `recall init`.

## Tasks

- [ ] 1. Add the `age` dependency and crypto package skeleton
  - [ ] 1.1 Add `filippo.io/age` to `go.mod` and create `internal/crypto` package
    - Define encrypt/decrypt function signatures using streaming I/O (`io.Reader`/`io.Writer`)
    - _Requirements: 1.1, 1.2, 1.3_

- [ ] 2. Implement application-managed identity
  - [ ] 2.1 Implement identity generation (`age` X25519 keypair)
    - _Requirements: 2.1_
  - [ ] 2.2 Define a credential-store interface with a default file-based fallback
    - Store public key in config; protect private key via OS credential store where available
    - _Requirements: 4.1, 4.2_
  - [ ] 2.3 Implement identity backup/export
    - _Requirements: 4.3, 4.4_

- [ ] 3. Wire identity generation to an explicit path (not `init`)
  - [ ] 3.1 Confirm `recall init` remains directory-only; do not generate identity in `init`
    - _Requirements: 3.1, 3.2_
  - [ ] 3.2 Generate identity via `recall keygen` and/or lazily on first `encrypt`, announcing creation
    - _Requirements: 3.3, 3.4_

- [ ] 4. Implement `recall encrypt`
  - [ ] 4.1 Application-key encryption to `<filename>.age` using safe (atomic) writes
    - _Requirements: 2.2, 2.4, 1.4, 6.1, 6.3_
  - [ ] 4.2 Passphrase encryption via `--passphrase` with secure prompt + confirmation
    - Never accept passphrase as an argument value; use `age` scrypt-based recipient
    - _Requirements: 5.1, 5.2, 5.4, 5.5, 6.2_

- [ ] 5. Implement `recall decrypt`
  - [ ] 5.1 Decrypt application-key files using the private identity
    - Error clearly if no identity exists
    - _Requirements: 2.3, 2.5, 3.5, 6.4_
  - [ ] 5.2 Decrypt passphrase files via secure prompt
    - _Requirements: 5.3, 6.4_

- [ ] 6. Integrate with existing commands
  - [ ] 6.1 Define and document `.age` handling in `list`
    - _Requirements: 7.1_
  - [ ] 6.2 Ensure `search` never matches ciphertext and never decrypts as a side effect
    - _Requirements: 7.2_
  - [ ] 6.3 Ensure recall/display never renders ciphertext as markdown
    - _Requirements: 7.3, 7.4_

- [ ] 7. Safety and hygiene
  - [ ] 7.1 Ensure private keys and passphrases are never logged or printed
    - _Requirements: 1.5_
  - [ ] 7.2 Ensure original files are preserved on encryption/decryption failure
    - _Requirements: 1.4_

- [ ] 8. Tests
  - [ ] 8.1 Unit tests for encrypt/decrypt round-trips (both models)
  - [ ] 8.2 Tests for `.age` interaction with list/search/display
  - [ ] 8.3 Tests confirming `init` does not create an identity

## Notes

- The encryption format should remain extensible for future multiple-recipient support, key rotation, and file sharing (Requirements 8.1–8.3).
- Credential-store backends (macOS Keychain, Linux Secret Service) can be added incrementally behind the interface from task 2.2.
- All crypto primitives come from `age`; no custom cryptography.
