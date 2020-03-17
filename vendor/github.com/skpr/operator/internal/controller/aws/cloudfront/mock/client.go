package mock

import (
	"errors"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/aws/aws-sdk-go/service/cloudfront/cloudfrontiface"
	"github.com/rs/xid"
)

// Client which mocks the CloudFront client.
type Client struct {
	cloudfrontiface.CloudFrontAPI
	distributions map[string]*cloudfront.Distribution
	tags          map[string]cloudfront.Tags
}

// New mock CloudFront client.
func New() *Client {
	return &Client{
		distributions: make(map[string]*cloudfront.Distribution),
		tags:          make(map[string]cloudfront.Tags),
	}
}

// CreateDistributionWithTags mocks the CloudFront client method.
func (m *Client) CreateDistributionWithTags(input *cloudfront.CreateDistributionWithTagsInput) (*cloudfront.CreateDistributionWithTagsOutput, error) {
	// Check if distribution exists.
	if _, ok := m.distributions[*input.DistributionConfigWithTags.DistributionConfig.CallerReference]; ok {
		return nil, awserr.New(cloudfront.ErrCodeDistributionAlreadyExists, "already exists", errors.New("distribution already exists"))
	}

	arn := aws.String(xid.New().String())

	distribution := &cloudfront.Distribution{
		ARN:                           arn,
		DistributionConfig:            input.DistributionConfigWithTags.DistributionConfig,
		Id:                            aws.String(xid.New().String()),
		InProgressInvalidationBatches: aws.Int64(1),
		DomainName:                    aws.String("xxxxxxx.cloudfront.net"),
		LastModifiedTime:              aws.Time(time.Now()),
		Status:                        aws.String("Deployed"),
	}

	m.distributions[*input.DistributionConfigWithTags.DistributionConfig.CallerReference] = distribution

	return &cloudfront.CreateDistributionWithTagsOutput{
		Distribution: distribution,
	}, nil
}

// GetDistributionConfig mocks the CloudFront client method.
func (m *Client) GetDistributionConfig(input *cloudfront.GetDistributionConfigInput) (*cloudfront.GetDistributionConfigOutput, error) {
	for _, distribution := range m.distributions {
		if *distribution.Id == *input.Id {
			return &cloudfront.GetDistributionConfigOutput{
				DistributionConfig: distribution.DistributionConfig,
			}, nil
		}
	}

	return nil, awserr.New(cloudfront.ErrCodeNoSuchDistribution, "not found", errors.New("distribution not found"))
}

// UpdateDistribution mocks the CloudFront client method.
func (m *Client) UpdateDistribution(input *cloudfront.UpdateDistributionInput) (*cloudfront.UpdateDistributionOutput, error) {
	for ref, distribution := range m.distributions {
		if *distribution.Id == *input.Id {
			distribution.DistributionConfig = input.DistributionConfig

			m.distributions[ref] = distribution

			return &cloudfront.UpdateDistributionOutput{
				Distribution: distribution,
			}, nil
		}
	}

	return nil, awserr.New(cloudfront.ErrCodeNoSuchDistribution, "not found", errors.New("distribution not found"))
}

// GetDistribution mocks the CloudFront client method.
func (m *Client) GetDistribution(input *cloudfront.GetDistributionInput) (*cloudfront.GetDistributionOutput, error) {
	for _, distribution := range m.distributions {
		if *distribution.Id == *input.Id {
			return &cloudfront.GetDistributionOutput{
				Distribution: distribution,
			}, nil
		}
	}

	return nil, awserr.New(cloudfront.ErrCodeNoSuchDistribution, "not found", errors.New("distribution not found"))
}

// ListDistributions mocks the CloudFront client method.
func (m *Client) ListDistributions(input *cloudfront.ListDistributionsInput) (*cloudfront.ListDistributionsOutput, error) {
	resp := &cloudfront.ListDistributionsOutput{
		DistributionList: &cloudfront.DistributionList{
			MaxItems:   aws.Int64(int64(len(m.distributions))),
			Quantity:   aws.Int64(int64(len(m.distributions))),
			NextMarker: aws.String(""),
		},
	}

	for _, distribution := range m.distributions {
		resp.DistributionList.Items = append(resp.DistributionList.Items, &cloudfront.DistributionSummary{
			Id: distribution.Id,
		})
	}

	return resp, nil
}

// DeleteDistribution mocks the CloudFront client method.
func (m *Client) DeleteDistribution(input *cloudfront.DeleteDistributionInput) (*cloudfront.DeleteDistributionOutput, error) {
	for ref, distribution := range m.distributions {
		if *distribution.Id == *input.Id {
			delete(m.distributions, ref)
		}
	}

	return nil, awserr.New(cloudfront.ErrCodeNoSuchDistribution, "not found", errors.New("distribution not found"))
}

// TagResource mocks the CloudFront client method.
func (m *Client) TagResource(*cloudfront.TagResourceInput) (*cloudfront.TagResourceOutput, error) {
	resp := &cloudfront.TagResourceOutput{}

	return resp, nil
}
