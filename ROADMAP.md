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

Support encrypting and decrypting individual recall files at rest, using a
password, an SSH key, or another mechanism already available on the system.
The goal is to protect sensitive notes without requiring external tooling.

#### File Encryption

The application should support local file encryption and decryption without requiring users to install or invoke external encryption tools.

The preferred implementation is to use the Go [`age`](https://pkg.go.dev/filippo.io/age) library internally. This provides a well-defined, modern file-encryption format while allowing the application to expose its own simple user experience.

Two encryption approaches are planned: **application-managed keys** and **passphrase-based encryption**.

#### Application-Managed Key Encryption

Application-managed key encryption should be the primary encryption method.

The application generates a dedicated public/private encryption identity for the user. These keys are separate from SSH keys and are used exclusively for protecting files managed by the application.

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

A typical workflow might look like:

```bash
app init
app encrypt example.txt
app decrypt example.txt.age
```

During initialization, the application generates a dedicated encryption identity. Subsequent encryption and decryption operations use that identity automatically.

##### Advantages

* Provides strong public-key encryption without requiring the user to manage passwords for each file.
* Separates file-encryption credentials from SSH authentication credentials.
* Allows encryption without access to the private key.
* Supports future workflows in which multiple users or systems can encrypt files for the same recipient.
* Provides a clean path toward supporting multiple recipients, key rotation, and shared encrypted files.
* Allows encryption and decryption to remain completely integrated into the application.

##### Key Storage

The private identity must be protected carefully.

Where practical, the application should store or protect the private key using operating-system credential storage rather than leaving the key unprotected in the application's configuration directory.

Possible platform-specific storage mechanisms include:

* macOS Keychain
* Linux Secret Service-compatible credential stores

The public key is not sensitive and may be stored in normal application configuration.

The application should also provide a documented backup or export mechanism for the private identity. Losing the private key would otherwise make encrypted files permanently unrecoverable.

#### Passphrase-Based Encryption

Passphrase encryption provides a simpler and more portable alternative for users who do not want to maintain an application identity.

Instead of using a public/private key pair, the user provides a passphrase when encrypting a file.

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

A typical workflow might look like:

```bash
app encrypt --passphrase example.txt
```

The application should securely prompt for the passphrase:

```text
Enter encryption passphrase:
Confirm encryption passphrase:
```

Decryption would similarly prompt for the passphrase:

```bash
app decrypt example.txt.age
```

```text
Enter decryption passphrase:
```

The passphrase should never be supplied as a normal command-line argument because command-line arguments may be exposed through shell history, process inspection, logs, or other system mechanisms.

The application should use the password-based encryption facilities provided by `age`, which use an appropriate password-based key derivation mechanism rather than treating the passphrase itself as an encryption key.

##### Advantages

* Requires no persistent encryption key.
* Makes encrypted files easy to move between systems.
* Provides a familiar user experience.
* Works well for files that need to be shared independently of an application identity.
* Provides a useful recovery or portability option alongside application-managed keys.

##### Limitations

Passphrase encryption depends heavily on the strength and availability of the user's passphrase.

A weak passphrase may reduce the effective security of the encrypted file, while a forgotten passphrase makes the file unrecoverable.

Unlike application-key encryption, passphrase encryption also requires access to the secret during both encryption and decryption.

#### Recommended Implementation

The application should support both approaches, with application-managed keys serving as the default.

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

The intended behavior is:

* **Application key encryption** — default option for files regularly managed by the application.
* **Passphrase encryption** — optional method for portability, sharing, or situations where persistent application keys are undesirable.

Both approaches should use the same underlying `age` file format where possible. This keeps the implementation consistent while allowing the application to expose different key-management models.

#### Design Goals

The encryption feature should follow these principles:

* Encryption and decryption should be built directly into the application.
* Users should not need to install the `age` command-line utility or other external encryption tools.
* Cryptographic formats and primitives should come from established libraries rather than custom implementations.
* Private keys and passphrases should never be written to logs.
* Passphrases should not be accepted as plain command-line arguments.
* Encryption should use streaming I/O so large files do not need to be loaded entirely into memory.
* Files should be written safely to avoid destroying the original file if encryption or decryption fails.
* The application should clearly communicate that losing either a private identity or passphrase may make encrypted data permanently unrecoverable.
* The encryption format should remain extensible enough to support future features such as multiple recipients, key rotation, and encrypted file sharing.


### Metrics collection

Track usage metrics for recall operations — files edited/created, recalled,
encrypted, and decrypted, among others — to give insight into how the tool
is used over time.

### TUI-based manager

Provide an interactive terminal UI to browse stored files and view collected
metrics, with the possibility of editing files directly from the interface.
