# Design Document

## Overview

This feature adds local file encryption and decryption to the Recall CLI without requiring users to install or invoke external encryption tools. The implementation uses the Go [`age`](https://pkg.go.dev/filippo.io/age) library internally, providing a well-defined, modern file-encryption format while allowing Recall to expose its own simple user experience.

Two encryption approaches are provided: **application-managed keys** (the default) and **passphrase-based encryption** (an optional, portable alternative). Both produce files in the `age` format.

## Application-Managed Key Encryption

Application-managed key encryption is the primary encryption method.

Recall generates a dedicated public/private encryption identity for the user. These keys are separate from SSH keys and are used exclusively for protecting files managed by Recall.

Conceptually:

```text
                    Application Identity
                           │
              ┌────────────┴────────────┐
              │                         │
         Public Key                Private Key
              │                         │
           Encrypt                    Decrypt
              │                         │
        Plaintext File            Encrypted File
              │                         │
              └────── Encrypted File ──┘
```

The public key can be used to encrypt data without exposing the private key. The private key is required only when decrypting files.

A typical workflow:

```bash
recall keygen           # or: identity created lazily on first encrypt
recall encrypt example
recall decrypt example.age
```

Recall generates a dedicated encryption identity through an explicit action (see "Identity Generation and `init`" below). Subsequent encryption and decryption operations use that identity automatically.

### Advantages

* Provides strong public-key encryption without requiring the user to manage passwords for each file.
* Separates file-encryption credentials from SSH authentication credentials.
* Allows encryption without access to the private key.
* Supports future workflows in which multiple users or systems can encrypt files for the same recipient.
* Provides a clean path toward supporting multiple recipients, key rotation, and shared encrypted files.
* Keeps encryption and decryption completely integrated into the application.

### Identity Generation and `init`

Recall's existing `recall init` command ensures the recall directory exists and is **not** required before normal usage. To avoid surprising behavior, identity generation is kept separate from `init`:

* `recall init` continues to only ensure the recall directory exists. It does **not** generate an encryption identity.
* The Application Identity is created through an explicit, documented path — either a dedicated `recall keygen` command, or lazily on the first `recall encrypt` when no identity exists.
* When an identity is created (either path), Recall clearly informs the user and reminds them to back it up.

The lazy-on-first-encrypt option is preferred for ergonomics, provided the creation is clearly announced; a dedicated `recall keygen` may be added for explicit control and re-generation/rotation workflows.

### Key Storage

The private identity must be protected carefully.

Where practical, Recall stores or protects the private key using operating-system credential storage rather than leaving the key unprotected in the recall directory or application configuration. Possible platform-specific mechanisms include:

* macOS Keychain
* Linux Secret Service-compatible credential stores

The public key is not sensitive and may be stored in normal application configuration.

Recall also provides a documented backup/export mechanism for the private identity. Losing the private key would otherwise make encrypted files permanently unrecoverable.

## Passphrase-Based Encryption

Passphrase encryption provides a simpler, more portable alternative for users who do not want to maintain an application identity.

Instead of a public/private key pair, the user provides a passphrase when encrypting a file.

Conceptually:

```text
Passphrase
    │
    ↓
Key Derivation
    │
    ↓
Encryption Key
    │
    ↓
Plaintext File → Encrypted File
```

The same passphrase is required to decrypt the file.

A typical workflow:

```bash
recall encrypt --passphrase example
```

Recall securely prompts for the passphrase:

```text
Enter encryption passphrase:
Confirm encryption passphrase:
```

Decryption similarly prompts:

```bash
recall decrypt example.age
```

```text
Enter decryption passphrase:
```

The passphrase is never supplied as a normal command-line argument, because command-line arguments may be exposed through shell history, process inspection, logs, or other system mechanisms. Recall uses the password-based encryption facilities provided by `age`, which apply an appropriate password-based key derivation mechanism rather than treating the passphrase itself as an encryption key.

### Advantages

* Requires no persistent encryption key.
* Makes encrypted files easy to move between systems.
* Provides a familiar user experience.
* Works well for files shared independently of an application identity.
* Provides a useful recovery or portability option alongside application-managed keys.

### Limitations

Passphrase encryption depends heavily on the strength and availability of the user's passphrase. A weak passphrase may reduce the effective security of the encrypted file, while a forgotten passphrase makes the file unrecoverable. Unlike application-key encryption, passphrase encryption requires access to the secret during both encryption and decryption.

## Interaction With Existing Commands

Encrypted files use a `.age` suffix. Because recall files otherwise have no extension and `storage.List` returns plain filenames, encrypted files would appear in `recall list` and be scanned by `recall search` unless handled explicitly. The design must specify:

* **`list`** — how `.age` files are displayed (e.g., shown with an indicator, or filtered).
* **`search`** — encrypted contents are ciphertext and MUST NOT be matched; searching MUST NOT decrypt files as a side effect.
* **recall/display** — displaying a `.age` file MUST NOT render raw ciphertext as markdown; Recall either decrypts (prompting as needed) or reports that the file is encrypted.

The chosen behavior for each interaction is documented for users.

## Recommended Implementation

Recall supports both approaches, with application-managed keys as the default.

```text
                         File Encryption
                               │
                  ┌────────────┴────────────┐
                  │                         │
          Application Identity         Passphrase
                  │                         │
           Public / Private            Key Derivation
                  │                         │
                  └────────────┬────────────┘
                               │
                            age API
                               │
                        Encrypted File
```

Intended behavior:

* **Application key encryption** — default option for files regularly managed by Recall.
* **Passphrase encryption** — optional method for portability, sharing, or situations where persistent application keys are undesirable.

Both approaches use the same underlying `age` file format where possible. This keeps the implementation consistent while allowing Recall to expose different key-management models.

## Design Goals

* Encryption and decryption are built directly into Recall.
* Users do not need to install the `age` command-line utility or other external encryption tools.
* Cryptographic formats and primitives come from established libraries (`age`) rather than custom implementations.
* Private keys and passphrases are never written to logs.
* Passphrases are never accepted as plain command-line arguments.
* Encryption uses streaming I/O so large files are not loaded entirely into memory.
* Files are written safely to avoid destroying the original file if encryption or decryption fails.
* Recall clearly communicates that losing either a private identity or passphrase may make encrypted data permanently unrecoverable.
* The encryption format remains extensible enough to support future features such as multiple recipients, key rotation, and encrypted file sharing.

## Proposed Package Layout

A new `internal/crypto` package encapsulates the `age` integration (identity generation, key storage abstraction, encrypt/decrypt streaming helpers). New `cmd/encrypt.go` and `cmd/decrypt.go` (and optionally `cmd/keygen.go`) wire the commands, following the existing Cobra patterns in `cmd/`. The credential-store abstraction is defined behind an interface so platform-specific backends (Keychain, Secret Service) can be added incrementally without changing callers.
