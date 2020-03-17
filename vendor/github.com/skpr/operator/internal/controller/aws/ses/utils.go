package ses

import "github.com/aws/aws-sdk-go/service/ses"

// Helper function to get the status of an AWS SES verification request.
func getVerificationStatus(resp *ses.GetIdentityVerificationAttributesOutput, address string) string {
	var status string

	if val, ok := resp.VerificationAttributes[address]; ok {
		status = *val.VerificationStatus
	}

	return status
}
