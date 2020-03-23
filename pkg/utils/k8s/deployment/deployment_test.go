// +build unit

package deployment

import (
	"testing"

	"github.com/stretchr/testify/assert"

	appsv1 "k8s.io/api/apps/v1"
)

func TestGetPhase(t *testing.T) {
	var replicas int32 = 2

	deployment := &appsv1.Deployment{
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
		},
		Status: appsv1.DeploymentStatus{
			Conditions: []appsv1.DeploymentCondition{
				{
					Type:   appsv1.DeploymentProgressing,
					Reason: TimedOutReason,
				},
			},
		},
	}
	assert.Equal(t, string(PhaseDeadLineExceeded), string(GetPhase(deployment)))

	deployment = &appsv1.Deployment{
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
		},
		Status: appsv1.DeploymentStatus{
			UpdatedReplicas: 1,
		},
	}
	assert.Equal(t, string(PhaseDeploying), string(GetPhase(deployment)))

	deployment = &appsv1.Deployment{
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
		},
		Status: appsv1.DeploymentStatus{
			UpdatedReplicas: 3,
		},
	}
	assert.Equal(t, string(PhaseDeploying), string(GetPhase(deployment)))

	deployment = &appsv1.Deployment{
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
		},
		Status: appsv1.DeploymentStatus{
			AvailableReplicas: 1,
			UpdatedReplicas:   3,
		},
	}
	assert.Equal(t, string(PhaseDeploying), string(GetPhase(deployment)))
}
