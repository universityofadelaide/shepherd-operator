package scheduled

import (
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	extensionsv1beta1 "github.com/skpr/operator/pkg/apis/extensions/v1beta1"
	skprmetav1 "github.com/skpr/operator/pkg/apis/meta/v1"
	resticutils "github.com/skpr/operator/pkg/utils/restic"
)

// Helper function to build a Backup with the required fields.
func buildBackup(scheduled *extensionsv1beta1.BackupScheduled, scheme *runtime.Scheme, scheduledTime time.Time) (*extensionsv1beta1.Backup, error) {
	// We want job names for a given nominal start time to have a deterministic name to avoid the same job being created twice
	name := fmt.Sprintf("%s-%d", scheduled.ObjectMeta.Name, scheduledTime.Unix())

	spec := *scheduled.Spec.Template.DeepCopy()
	var tags []string
	if spec.Tags != nil {
		tags = spec.Tags
	}
	spec.Tags = append(tags, resticutils.ScheduledTag)
	backup := &extensionsv1beta1.Backup{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Backup",
			APIVersion: "extensions.skpr.io/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: scheduled.ObjectMeta.Namespace,
			Annotations: map[string]string{
				skprmetav1.ScheduledAnnotation: scheduledTime.Format(time.RFC3339),
			},
		},
		Spec: spec,
	}

	if err := controllerutil.SetControllerReference(scheduled, backup, scheme); err != nil {
		return nil, err
	}

	return backup, nil
}
