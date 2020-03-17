package restic

// Prefix for discoverying Restic resources.
const Prefix = "backup-restic"

// RestorePrefix for discoverying Restic Restore resources.
const RestorePrefix = "restore-restic"

// ResticBackupContainerName is the name of the container in the restic backup pod.
const ResticBackupContainerName = "restic-backup"

// ScheduledTag is the tag to applied to scheduled backups.
const ScheduledTag = "scheduled"
