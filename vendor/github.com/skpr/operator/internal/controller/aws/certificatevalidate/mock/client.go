package mock

import (
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/route53/route53iface"
)

// Client which mocks the CloudFront client.
type Client struct {
	route53iface.Route53API
	Records map[string][]string
}

// New mock CloudFront client.
func New() *Client {
	return &Client{
		Records: make(map[string][]string),
	}
}

// ChangeResourceRecordSets mock implementation.
func (c *Client) ChangeResourceRecordSets(input *route53.ChangeResourceRecordSetsInput) (*route53.ChangeResourceRecordSetsOutput, error) {
	var resp *route53.ChangeResourceRecordSetsOutput

	for _, change := range input.ChangeBatch.Changes {
		for _, record := range change.ResourceRecordSet.ResourceRecords {
			c.Records[*change.ResourceRecordSet.Name] = append(c.Records[*change.ResourceRecordSet.Name], *record.Value)
		}
	}

	return resp, nil
}
