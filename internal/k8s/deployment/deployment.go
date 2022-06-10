package deployment

import (
	appsv1 "k8s.io/api/apps/v1"
)

// Phase of a deployment.
type Phase string

const (
	// TimedOutReason is used to determine if the deployment is PhaseDeadLineExceeded.
	TimedOutReason = "ProgressDeadlineExceeded"
	// PhaseDeadLineExceeded declares that a deployment has exceeded the deadline.
	PhaseDeadLineExceeded Phase = "DeadLineExceeded"
	// PhaseDeploying declares that a deploment is still has deploying.
	PhaseDeploying Phase = "Deploying"
	// PhaseDeployed declares that a deployment has finished deploying.
	PhaseDeployed Phase = "Deployed"
)

// GetPhase returns the current state of a deployment.
func GetPhase(deployment *appsv1.Deployment) Phase {
	cond := getDeploymentCondition(deployment.Status, appsv1.DeploymentProgressing)

	if cond != nil && cond.Reason == TimedOutReason {
		return PhaseDeadLineExceeded
	}

	if deployment.ObjectMeta.Generation > deployment.Status.ObservedGeneration {
		return PhaseDeploying
	}

	if deployment.Status.UpdatedReplicas < *deployment.Spec.Replicas {
		return PhaseDeploying
	}

	if deployment.Status.Replicas > deployment.Status.UpdatedReplicas {
		return PhaseDeploying
	}

	if deployment.Status.AvailableReplicas < deployment.Status.UpdatedReplicas {
		return PhaseDeploying
	}

	return PhaseDeployed
}

// GetDeploymentConditionInternal returns the condition with the provided type.
// Borrowed from: https://github.com/kubernetes/kubernetes/blob/master/pkg/controller/deployment/util/deployment_util.go#L135
func getDeploymentCondition(status appsv1.DeploymentStatus, condType appsv1.DeploymentConditionType) *appsv1.DeploymentCondition {
	for i := range status.Conditions {
		c := status.Conditions[i]

		if c.Type == condType {
			return &c
		}
	}

	return nil
}
