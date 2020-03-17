package job

import (
	"testing"

	skprmetav1 "github.com/skpr/operator/pkg/apis/meta/v1"
	"github.com/stretchr/testify/assert"
	batchv1 "k8s.io/api/batch/v1"
)

func TestGetPhase(t *testing.T) {
	assert.Equal(t, skprmetav1.PhaseInProgress, GetPhase(batchv1.JobStatus{Active: 1}))
	assert.Equal(t, skprmetav1.PhaseFailed, GetPhase(batchv1.JobStatus{Failed: 1}))
	assert.Equal(t, skprmetav1.PhaseCompleted, GetPhase(batchv1.JobStatus{Succeeded: 1}))
	assert.Equal(t, skprmetav1.PhaseUnknown, GetPhase(batchv1.JobStatus{}))
}
