package pod

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	shpdmetav1 "github.com/universityofadelaide/shepherd-operator/apis/meta/v1"
)

// CompletionTime for a Pod using container status.
func CompletionTime(pod *corev1.Pod) *metav1.Time {
	oldest := pod.ObjectMeta.CreationTimestamp

	for _, container := range pod.Status.ContainerStatuses {
		if container.State.Terminated == nil {
			continue
		}

		if container.State.Terminated.FinishedAt.Before(&oldest) {
			continue
		}

		oldest = container.State.Terminated.FinishedAt
	}

	return &oldest
}

// GetPhase from the Pod status object.
func GetPhase(status corev1.PodStatus) shpdmetav1.Phase {
	switch status.Phase {
	case corev1.PodSucceeded:
		return shpdmetav1.PhaseCompleted
	case corev1.PodFailed:
		return shpdmetav1.PhaseFailed
	default:
		return shpdmetav1.PhaseInProgress
	}
}
