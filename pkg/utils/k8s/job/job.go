package job

import (
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	skprmetav1 "github.com/skpr/operator/pkg/apis/meta/v1"
)

// NewFromPod enforces Job best practices.
func NewFromPod(metadata metav1.ObjectMeta, spec corev1.PodSpec) (*batchv1.Job, error) {
	var (
		parallelism int32 = 1
		completions int32 = 1
		deadline    int64 = 1800
		backoff     int32 = 2
	)

	job := &batchv1.Job{
		ObjectMeta: metadata,
		Spec: batchv1.JobSpec{
			Parallelism:           &parallelism,
			Completions:           &completions,
			ActiveDeadlineSeconds: &deadline,
			BackoffLimit:          &backoff,
			Template: corev1.PodTemplateSpec{
				Spec: spec,
			},
		},
	}

	return job, nil
}

// GetPhase from the JobStatus object.
func GetPhase(status batchv1.JobStatus) skprmetav1.Phase {
	if status.Active > 0 {
		return skprmetav1.PhaseInProgress
	}

	if status.Failed > 0 {
		return skprmetav1.PhaseFailed
	}

	if status.Succeeded > 0 {
		return skprmetav1.PhaseCompleted
	}

	return skprmetav1.PhaseUnknown
}
