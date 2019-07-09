package restic

// Prefix for discoverying Restic resources.
const Prefix = "backup-restic"

// VolumeSecrets identifier used for Restic secret.
const VolumeSecrets = "restic-secrets"

// VolumeRepository identifier used for Restic repository.
const VolumeRepository = "restic-repository"

// ResticSecretPasswordName is the name of the secret the restic password is stored in.
const ResticSecretPasswordName = "shepherd-restic-secret"

// ResticBackupContainerName is the name of the container in the restic backup pod.
const ResticBackupContainerName = "restic-backup"
