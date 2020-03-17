package imagescheduled

import (
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	skprmetav1 "github.com/skpr/operator/pkg/apis/meta/v1"
	mysqlv1beta1 "github.com/skpr/operator/pkg/apis/mysql/v1beta1"
)

// Helper function to build an Image with the required fields.
func buildImage(scheduled *mysqlv1beta1.ImageScheduled, scheme *runtime.Scheme, scheduledTime time.Time) (*mysqlv1beta1.Image, error) {
	// We want job names for a given nominal start time to have a deterministic name to avoid the same job being created twice
	name := fmt.Sprintf("%s-%d", scheduled.ObjectMeta.Name, scheduledTime.Unix())

	image := &mysqlv1beta1.Image{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Image",
			APIVersion: "mysql.skpr.io/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: scheduled.ObjectMeta.Namespace,
			Annotations: map[string]string{
				skprmetav1.ScheduledAnnotation: scheduledTime.Format(time.RFC3339),
			},
		},
		Spec: *scheduled.Spec.Template.DeepCopy(),
	}

	if err := controllerutil.SetControllerReference(scheduled, image, scheme); err != nil {
		return nil, err
	}

	return image, nil
}
