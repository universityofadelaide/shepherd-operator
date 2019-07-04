package restic

// Prefix for discoverying Restic resources.
const Prefix = "backup-restic"

// VolumeSecrets identifier used for Restic secret.
const VolumeSecrets = "restic-secrets"

// ResticSecretPasswordName is the name of the secret the restic password is stored in.
const ResticSecretPasswordName = "shepherd-restic-secret"
