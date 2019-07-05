package iam

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
)

// Client which mocks the IAM client.
type Client struct {
	iamiface.IAMAPI
}

// New mock IAM client.
func New() *Client {
	return &Client{}
}

// CreateUser is a mock implementation.
func (c *Client) CreateUser(input *iam.CreateUserInput) (*iam.CreateUserOutput, error) {
	return &iam.CreateUserOutput{
		User: &iam.User{
			UserName: input.UserName,
		},
	}, nil
}

// ListAccessKeys is a mock implementation.
func (c *Client) ListAccessKeys(*iam.ListAccessKeysInput) (*iam.ListAccessKeysOutput, error) {
	return &iam.ListAccessKeysOutput{}, nil
}

// CreateAccessKey is a mock implementation.
func (c *Client) CreateAccessKey(input *iam.CreateAccessKeyInput) (*iam.CreateAccessKeyOutput, error) {
	return &iam.CreateAccessKeyOutput{
		AccessKey: &iam.AccessKey{
			AccessKeyId:     aws.String("xxxxxxxxxxxxxxxxxx"),
			SecretAccessKey: aws.String("yyyyyyyyyyyyyyyyyy"),
		},
	}, nil
}

// DeleteAccessKey is a mock implementation.
func (c *Client) DeleteAccessKey(*iam.DeleteAccessKeyInput) (*iam.DeleteAccessKeyOutput, error) {
	return &iam.DeleteAccessKeyOutput{}, nil
}

// PutUserPolicy is a mock implementation.
func (c *Client) PutUserPolicy(*iam.PutUserPolicyInput) (*iam.PutUserPolicyOutput, error) {
	return &iam.PutUserPolicyOutput{}, nil
}
