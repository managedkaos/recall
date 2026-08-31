# Requirements Document

## Introduction

This feature adds local file encryption and decryption to the Recall CLI. Users can encrypt individual recall files at rest and decrypt them on demand, without installing or invoking any external encryption tools. Encryption is built directly into the application using the Go [`age`](https://pkg.go.dev/filippo.io/age) library, which provides a modern, well-defined file-encryption format while allowing Recall to expose its own simple user experience.

Two encryption models are supported: **application-managed key encryption** (the default) and **passphrase-based encryption** (an optional, portable alternative). Both produce files in the `age` format.

## Glossary

- **Recall_CLI**: The command-line application that provides all recall functionality
- **Recall_File**: A markdown-formatted file stored in the Recall_Directory without a `.md` extension
- **Recall_Directory**: The directory on the user's filesystem where recall files are stored (default `~/.recall` or `RECALL_DIR`)
- **Encrypted_File**: A Recall file encrypted with the `age` format, identified by a `.age` suffix
- **Age_Library**: The `filippo.io/age` Go library used for all cryptographic operations
- **Application_Identity**: A dedicated `age` public/private key pair generated and managed by Recall, distinct from SSH keys
- **Public_Key**: The recipient key used to encrypt files; not sensitive
- **Private_Key**: The identity key required to decrypt files; sensitive and protected
- **Passphrase**: A user-supplied secret used to derive an encryption key for passphrase-based encryption
- **Credential_Store**: An operating-system credential storage mechanism (e.g., macOS Keychain, Linux Secret Service)

## Requirements

### Requirement 1: Built-in Encryption

**User Story:** As a user, I want encryption built directly into Recall, so that I do not have to install or invoke external encryption tools.

#### Acceptance Criteria

1. THE Recall_CLI SHALL perform all encryption and decryption using the Age_Library internally, without shelling out to the `age` command-line utility or any other external tool.
2. THE Recall_CLI SHALL use cryptographic formats and primitives provided by the Age_Library rather than custom implementations.
3. THE Recall_CLI SHALL use streaming I/O for encryption and decryption so that large files are not loaded entirely into memory.
4. WHEN writing an Encrypted_File or a decrypted output file, THE Recall_CLI SHALL write safely so that a failure during encryption or decryption does not destroy or corrupt the original input file.
5. THE Recall_CLI SHALL NOT write Private_Keys or Passphrases to logs or standard output.

### Requirement 2: Application-Managed Key Encryption

**User Story:** As a user, I want Recall to manage a dedicated encryption identity, so that I can encrypt and decrypt files without managing a password per file.

#### Acceptance Criteria

1. THE Recall_CLI SHALL support generating an Application_Identity consisting of an `age` Public_Key and Private_Key, distinct from any SSH keys.
2. WHEN the user runs `recall encrypt <filename>` and an Application_Identity exists, THE Recall_CLI SHALL encrypt the file to the identity's Public_Key and write the result as `<filename>.age`.
3. WHEN the user runs `recall decrypt <filename>.age` and an Application_Identity exists, THE Recall_CLI SHALL decrypt the file using the identity's Private_Key.
4. THE Recall_CLI SHALL be able to encrypt files using only the Public_Key, without requiring access to the Private_Key.
5. THE Recall_CLI SHALL require the Private_Key only for decryption.

### Requirement 3: Identity Generation and Its Relationship to `init`

**User Story:** As a user, I want a clear and predictable relationship between `recall init` and encryption identity generation, so that initializing the recall directory does not silently create cryptographic material I did not expect.

#### Acceptance Criteria

1. THE Recall_CLI SHALL continue to treat `recall init` as an operation that ensures the Recall_Directory exists, and `recall init` SHALL NOT be required before normal (non-encryption) usage.
2. THE Recall_CLI SHALL NOT generate an Application_Identity as an implicit side effect of `recall init`.
3. THE Recall_CLI SHALL generate the Application_Identity through an explicit, documented action — either a dedicated command (e.g., `recall keygen`) OR lazily on the first `recall encrypt` invocation when no identity exists — and SHALL clearly communicate when an identity is created.
4. IF an identity is generated lazily on first encrypt, THEN THE Recall_CLI SHALL inform the user that a new Application_Identity was created and remind them to back it up.
5. IF the user runs `recall decrypt` and no Application_Identity exists, THEN THE Recall_CLI SHALL exit with a non-zero code and an error message indicating that no encryption identity is available.

### Requirement 4: Private Key Storage and Backup

**User Story:** As a user, I want my private encryption key protected and recoverable, so that my encrypted files stay both secure and recoverable.

#### Acceptance Criteria

1. WHERE the operating system provides a Credential_Store, THE Recall_CLI SHOULD store or protect the Private_Key using that Credential_Store rather than leaving it unprotected in the Recall_Directory or application configuration.
2. THE Recall_CLI MAY store the Public_Key in normal application configuration, as the Public_Key is not sensitive.
3. THE Recall_CLI SHALL provide a documented mechanism to back up or export the Application_Identity's Private_Key.
4. THE Recall_CLI SHALL clearly communicate that losing the Private_Key makes files encrypted to that identity permanently unrecoverable.

### Requirement 5: Passphrase-Based Encryption

**User Story:** As a user, I want to encrypt a file with a passphrase, so that I can protect and move files without maintaining an application identity.

#### Acceptance Criteria

1. WHEN the user runs `recall encrypt --passphrase <filename>`, THE Recall_CLI SHALL derive an encryption key from a user-supplied Passphrase using the Age_Library's password-based facilities and write the result as `<filename>.age`.
2. WHEN encrypting with a Passphrase, THE Recall_CLI SHALL securely prompt for the Passphrase and prompt again for confirmation, without echoing the Passphrase to the terminal.
3. WHEN the user runs `recall decrypt <filename>.age` on a passphrase-encrypted file, THE Recall_CLI SHALL securely prompt for the Passphrase without echoing it.
4. THE Recall_CLI SHALL NOT accept a Passphrase as a plain command-line argument value.
5. THE Recall_CLI SHALL rely on the Age_Library's password-based key derivation rather than treating the raw Passphrase as an encryption key.

### Requirement 6: Encryption Model Selection

**User Story:** As a user, I want application-key encryption by default with passphrase encryption available, so that I get strong protection by default and portability when I need it.

#### Acceptance Criteria

1. WHEN the user runs `recall encrypt <filename>` without a passphrase option, THE Recall_CLI SHALL use application-managed key encryption as the default.
2. WHEN the user runs `recall encrypt --passphrase <filename>`, THE Recall_CLI SHALL use passphrase-based encryption instead of application-managed key encryption.
3. THE Recall_CLI SHALL produce Encrypted_Files in the `age` format for both encryption models.
4. WHEN decrypting, THE Recall_CLI SHALL determine the appropriate decryption path (identity vs. passphrase) in a manner consistent with the `age` format of the Encrypted_File.

### Requirement 7: Interaction With Existing Commands

**User Story:** As a user, I want encrypted files to behave predictably with `list`, `search`, and recall/display, so that encrypted content is not accidentally exposed or misrendered.

#### Acceptance Criteria

1. THE Recall_CLI SHALL define and document how Encrypted_Files (`.age` suffix) appear in `recall list` output.
2. WHEN the user runs `recall search`, THE Recall_CLI SHALL NOT match against the ciphertext contents of Encrypted_Files, and SHALL NOT decrypt files as a side effect of searching.
3. WHEN the user attempts to display an Encrypted_File via the recall (display) operation, THE Recall_CLI SHALL either decrypt it (prompting as needed) OR report that the file is encrypted, and SHALL NOT render raw ciphertext as markdown.
4. THE Recall_CLI SHALL document the chosen behavior for each of the above interactions.

### Requirement 8: Extensibility

**User Story:** As a maintainer, I want the encryption design to remain extensible, so that future features can build on it without rework.

#### Acceptance Criteria

1. THE encryption format SHALL remain extensible enough to support future multiple-recipient encryption.
2. THE design SHALL allow for future key rotation of the Application_Identity.
3. THE design SHALL allow for future encrypted file sharing between users or systems.
