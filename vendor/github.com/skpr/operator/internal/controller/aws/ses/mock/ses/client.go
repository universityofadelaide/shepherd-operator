package ses

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/aws/aws-sdk-go/service/ses/sesiface"
)

// Client which mocks the SES client.
type Client struct {
	sesiface.SESAPI
	attributes map[string]*ses.IdentityVerificationAttributes
}

// New mock SES client.
func New() *Client {
	return &Client{
		attributes: make(map[string]*ses.IdentityVerificationAttributes),
	}
}

// GetIdentityVerificationAttributes is a mock implementation.
func (c *Client) GetIdentityVerificationAttributes(input *ses.GetIdentityVerificationAttributesInput) (*ses.GetIdentityVerificationAttributesOutput, error) {
	return &ses.GetIdentityVerificationAttributesOutput{
		VerificationAttributes: c.attributes,
	}, nil
}

// VerifyEmailAddress is a mock implementation.
func (c *Client) VerifyEmailAddress(input *ses.VerifyEmailAddressInput) (*ses.VerifyEmailAddressOutput, error) {
	c.attributes[*input.EmailAddress] = &ses.IdentityVerificationAttributes{
		VerificationStatus: aws.String("Pending"),
	}

	return &ses.VerifyEmailAddressOutput{}, nil
}
