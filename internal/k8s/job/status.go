package job

import (
	batchv1 "k8s.io/api/batch/v1"
)

// IsFinished which check if a Job has finished (completed or failed).
func IsFinished(job *batchv1.Job) bool {
	for _, condition := range job.Status.Conditions {
		switch condition.Type {
		case batchv1.JobComplete:
			return true
		case batchv1.JobFailed:
			return true
		}
	}

	return false
}
