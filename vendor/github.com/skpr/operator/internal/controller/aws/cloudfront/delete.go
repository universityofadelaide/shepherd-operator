package cloudfront

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	corev1 "k8s.io/api/core/v1"

	awsv1beta1 "github.com/skpr/operator/pkg/apis/aws/v1beta1"
	"github.com/skpr/operator/pkg/utils/k8s/events"
)

// DeleteExternal CloudFront distribution.
func (r *ReconcileCloudFront) DeleteExternal(instance *awsv1beta1.CloudFront) error {
	r.recorder.Eventf(instance, corev1.EventTypeNormal, events.EventDelete, "Deleting distribution: %s", instance.Status.ID)

	_, err := r.cloudfront.DeleteDistribution(&cloudfront.DeleteDistributionInput{
		Id: aws.String(instance.Status.ID),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() == cloudfront.ErrCodeNoSuchResource {
				return nil
			}
		} else {
			return err
		}
	}

	return nil
}
