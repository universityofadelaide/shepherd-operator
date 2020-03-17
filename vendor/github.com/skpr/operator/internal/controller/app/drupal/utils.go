package drupal

import (
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// If the user has asked for 1 copy of the application to be running we want to ensure the
// old copy of the app does not go away until the new is established.
func applyStrategy(deployment *appsv1.Deployment) error {
	if *deployment.Spec.Replicas < int32(2) {
		var (
			maxUnavailable = intstr.FromInt(0)
			maxSurge       = intstr.FromInt(1)
		)

		deployment.Spec.Strategy.RollingUpdate = &appsv1.RollingUpdateDeployment{
			MaxUnavailable: &maxUnavailable,
			MaxSurge:       &maxSurge,
		}
	}

	return nil
}
