package cloudfront

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/aws/aws-sdk-go/service/cloudfront/cloudfrontiface"
	"github.com/go-test/deep"
	"github.com/pkg/errors"
	"github.com/prometheus/common/log"
	corev1 "k8s.io/api/core/v1"

	awsv1beta1 "github.com/skpr/operator/pkg/apis/aws/v1beta1"
	"github.com/skpr/operator/pkg/utils/k8s/events"
)

// SyncExternal CloudFront distribution.
func (r *ReconcileCloudFront) SyncExternal(log log.Logger, instance *awsv1beta1.CloudFront) (*cloudfront.Distribution, error) {
	config, err := generateDistribution(r.prefix, instance)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate distribution config")
	}

	if instance.Status.ID != "" {
		log.Info("ID found on object status field. Updating existing distribution.")
		return r.UpdateDistribution(instance, config)
	}

	// The distribution does not exist, we should create it.
	resp, err := r.cloudfront.CreateDistribution(&cloudfront.CreateDistributionInput{
		DistributionConfig: config,
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() == cloudfront.ErrCodeDistributionAlreadyExists {
				log.Info("CreateDistribution failed. Distribution already exists. Finding ID by CallerReference.")

				id, err := findDistributionIDFromRef(r.cloudfront, *config.CallerReference)
				if err != nil {
					return nil, errors.Wrap(err, "failed to get reference")
				}

				instance.Status.ID = id

				return r.UpdateDistribution(instance, config)
			}
		}

		return nil, errors.Wrap(err, "failed to create distribution")
	}

	return resp.Distribution, nil
}

// UpdateDistribution CloudFront distribution with new DistributionConfig.
func (r *ReconcileCloudFront) UpdateDistribution(instance *awsv1beta1.CloudFront, config *cloudfront.DistributionConfig) (*cloudfront.Distribution, error) {
	// Get the current Distribution configuration.
	respGet, err := r.cloudfront.GetDistributionConfig(&cloudfront.GetDistributionConfigInput{
		Id: aws.String(instance.Status.ID),
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get distribution config")
	}

	// Compare it to what we have. Update if nessecary.
	if diff := deep.Equal(respGet.DistributionConfig, config); diff != nil {
		log.Info(fmt.Sprintf("Distribution change dectected: %s", diff))

		r.recorder.Eventf(instance, corev1.EventTypeNormal, events.EventUpdate, "Updating distribution: %s", instance.Status.ID)

		update, err := r.cloudfront.UpdateDistribution(&cloudfront.UpdateDistributionInput{
			Id:                 aws.String(instance.Status.ID),
			IfMatch:            respGet.ETag,
			DistributionConfig: config,
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to update distribution")
		}

		return update.Distribution, nil
	}

	// Get the current state of the distribution.
	resp, err := r.cloudfront.GetDistribution(&cloudfront.GetDistributionInput{
		Id: aws.String(instance.Status.ID),
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to update distribution")
	}

	return resp.Distribution, nil
}

// This is an expensive function which is only used for finding a CloudFront via CallerReference if it has been lost.
//   eg. Distribution created by .Status.ARN not updated.
func findDistributionIDFromRef(client cloudfrontiface.CloudFrontAPI, ref string) (string, error) {
	input := &cloudfront.ListDistributionsInput{}

	for {
		resp, err := client.ListDistributions(input)
		if err != nil {
			return "", errors.Wrap(err, "failed to list distributions")
		}

		for _, dist := range resp.DistributionList.Items {
			config, err := client.GetDistributionConfig(&cloudfront.GetDistributionConfigInput{
				Id: dist.Id,
			})
			if err != nil {
				return "", errors.Wrap(err, "failed to get distribution config")
			}

			if *config.DistributionConfig.CallerReference == ref {
				return *dist.Id, nil
			}
		}

		if resp.DistributionList.NextMarker == nil {
			break
		}

		input.Marker = resp.DistributionList.NextMarker
	}

	return "", fmt.Errorf("distribution config not found: %s", ref)
}
