package restic

// Prefix for discovering Restic resources.
const Prefix = "restic"

// VolumeSecrets identifier used for Restic secret.
const VolumeSecrets = "restic-secrets"

// VolumeRepository identifier used for Restic repository.
const VolumeRepository = "restic-repository"

// ResticSecretPasswordName is the name of the secret the restic password is stored in.
const ResticSecretPasswordName = "shepherd-restic-secret"

// ResticBackupContainerName is the name of the container in the restic backup pod.
const ResticBackupContainerName = "restic-backup"

// WebDirectory is working directory for the restore deployment step.
const WebDirectory = "/code"

// FriendlyNameAnnotation is the name of the annotation which stores the friendly name of a backup to display in the Shepherd UI.
const FriendlyNameAnnotation = "backups.shepherd/friendly-name"

// SyncLabel is the name of the label to determine that a backup is part of a sync.
const SyncLabel = "is-sync"
