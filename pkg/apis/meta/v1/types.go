package v1

type Phase string

const (
	PhaseNew        Phase = "new"
	PhaseInProgress Phase = "InProgress"
	PhaseFailed     Phase = "Failed"
	PhaseCompleted  Phase = "Completed"
)

// SpecVolume defines how to  volumes.
type SpecVolume struct {
	// ClaimName which will be backed up.
	ClaimName string `json:"claimName"`
}

// SpecMySQL defines how to  MySQL.
type SpecMySQL struct {
	// Secret which will be used for connectivity.
	Secret SpecMySQLSecret `json:"secret"`
}

type SpecMySQLSecret struct {
	// Name of secret containing the mysql connection details.
	Name string `json:"name"`
	// Keys within secret to use for each parameter.
	Keys SpecMySQLSecretKeys `json:"keys"`
}

// SpecMySQLSecretKeys defines Secret keys for MySQL connectivity.
type SpecMySQLSecretKeys struct {
	// Key which was applied to the application for database connectivity.
	Username string `json:"username"`
	// Key which was applied to the application for database connectivity.
	Password string `json:"password"`
	// Key which was applied to the application for database connectivity.
	Database string `json:"database"`
	// Key which was applied to the application for database connectivity.
	Hostname string `json:"hostname"`
	// Key which was applied to the application for database connectivity.
	Port string `json:"port"`
}
