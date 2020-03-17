package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ScheduledAnnotation is used to detect the time when an object was scheduled.
const ScheduledAnnotation = "skpr.io/scheduled-at"

// Phase which indicates the status of an operation.
type Phase string

const (
	// PhaseFailed to be assigned when an operation fails.
	PhaseFailed Phase = "Failed"
	// PhaseReady to be assigned when an operation is ready to be progressed.
	PhaseReady Phase = "Ready"
	// PhaseInProgress to be assigned when an operation is in progress.
	PhaseInProgress Phase = "InProgress"
	// PhaseCompleted to be assigned when an operation has been completed.
	PhaseCompleted Phase = "Completed"
	// PhaseUnknown to be assigned when the above phases cannot be determined.
	PhaseUnknown Phase = "Unknown"
)

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
)

// +genclient
// +k8s:deepcopy-gen=true

// ScheduledSpec defines the desired state of a scheduled task.
type ScheduledSpec struct {
	// The schedule in Cron format, see https://en.wikipedia.org/wiki/Cron.
	CronTab string `json:"cronTab"`

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
// +k8s:deepcopy-gen=true

// ScheduledStatus defines the observed state of a scheduled object.
type ScheduledStatus struct {
	// A list of pointers to currently running jobs.
	// +optional
	Active []corev1.ObjectReference `json:"active,omitempty"`

	// Information when was the last time the job was successfully scheduled.
	// +optional
	LastScheduleTime *metav1.Time `json:"lastScheduleTime,omitempty"`
}
