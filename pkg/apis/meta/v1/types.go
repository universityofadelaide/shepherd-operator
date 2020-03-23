package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Phase string

const (
	PhaseNew        Phase = "New"
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

// ConcurrencyPolicy describes how the scheduled task will be handled.
// Only one of the following concurrent policies may be specified.
// If none of the following policies is specified, the default one
// is ForbidConcurrent.
type ConcurrencyPolicy string

const (
	// AllowConcurrent allows CronJobs to run concurrently.
	AllowConcurrent ConcurrencyPolicy = "Allow"

	// ForbidConcurrent forbids concurrent runs, skipping next run if previous
	// hasn't finished yet.
	ForbidConcurrent ConcurrencyPolicy = "Forbid"

	// ReplaceConcurrent cancels currently running job and replaces it with a new one.
	ReplaceConcurrent ConcurrencyPolicy = "Replace"

	// FriendlyNameFormat is the time format for default friendly name.
	FriendlyNameFormat = "Mon, 02/01/2006 - 15:04"
)

// ScheduledAnnotation is used to detect the time when an object was scheduled.
const ScheduledAnnotation = "skpr.io/scheduled-at"

// +genclient
// +k8s:deepcopy-gen=true

// ScheduledSpec defines the desired state of a scheduled task.
type ScheduledSpec struct {
	// The schedule in Cron format, see https://en.wikipedia.org/wiki/Cron.
	CronTab string `json:"crontab"`

	// Optional deadline in seconds for starting the job if it misses scheduled
	// time for any reason.  Missed jobs executions will be counted as failed ones.
	StartingDeadlineSeconds *int64 `json:"startingDeadlineSeconds,omitempty"`

	// Specifies how to treat concurrent executions of a Job.
	// Valid values are:
	// - "Allow" (default): allows CronJobs to run concurrently;
	// - "Forbid": forbids concurrent runs, skipping next run if previous run hasn't finished yet;
	// - "Replace": cancels currently running job and replaces it with a new one
	ConcurrencyPolicy ConcurrencyPolicy `json:"concurrencyPolicy,omitempty"`

	// This flag tells the controller to suspend subsequent executions, it does
	// not apply to already started executions.  Defaults to false.
	Suspend *bool `json:"suspend,omitempty"`

	// The number of successful finished jobs to retain.
	// This is a pointer to distinguish between explicit zero and not specified.
	SuccessfulHistoryLimit *int32 `json:"successfulJobsHistoryLimit,omitempty"`

	// The number of failed finished jobs to retain.
	// This is a pointer to distinguish between explicit zero and not specified.
	FailedHistoryLimit *int32 `json:"failedJobsHistoryLimit,omitempty"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ScheduledStatus defines the observed state of a scheduled object.
type ScheduledStatus struct {
	// A list of pointers to currently running jobs.
	// +optional
	Active []corev1.ObjectReference `json:"active,omitempty"`

	// Information when was the last time the job was successfully scheduled.
	// +optional
	LastExecutedTime *metav1.Time `json:"lastExecutedTime,omitempty"`
}
