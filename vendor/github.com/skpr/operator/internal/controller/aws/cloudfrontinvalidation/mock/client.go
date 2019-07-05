package mock

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/aws/aws-sdk-go/service/cloudfront/cloudfrontiface"
	"github.com/rs/xid"

	awsv1beta1 "github.com/skpr/operator/pkg/apis/aws/v1beta1"
)

// Client which mocks the CloudFront client.
type Client struct {
	cloudfrontiface.CloudFrontAPI
	invalidations map[string]*cloudfront.Invalidation
}

// New mock CloudFront client.
func New() *Client {
	return &Client{
		invalidations: make(map[string]*cloudfront.Invalidation),
	}
}

// CreateInvalidation mocks the CreateInvalidation function.
func (m *Client) CreateInvalidation(input *cloudfront.CreateInvalidationInput) (*cloudfront.CreateInvalidationOutput, error) {
	invalidation := &cloudfront.Invalidation{
		Id:         aws.String(xid.New().String()),
		CreateTime: aws.Time(time.Now()),
		Status:     aws.String("InProgress"),
	}

	m.invalidations[*invalidation.Id] = invalidation

	return &cloudfront.CreateInvalidationOutput{
		Invalidation: invalidation,
	}, nil
}

// GetInvalidation mocks the GetInvalidation function.
func (m *Client) GetInvalidation(input *cloudfront.GetInvalidationInput) (*cloudfront.GetInvalidationOutput, error) {
	for id, invalidation := range m.invalidations {
		if id == *input.Id {
			resp := &cloudfront.GetInvalidationOutput{
				Invalidation: invalidation,
			}

			resp.Invalidation.Status = aws.String(awsv1beta1.CloudFrontInvalidationCompleted)

			return resp, nil
		}
	}

	return nil, fmt.Errorf("invalidation not found")
}
