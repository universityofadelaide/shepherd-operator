package certificate

import (
	"sort"

	"github.com/aws/aws-sdk-go/service/acm"
	"github.com/maple-tech/go-hashify"
	awsv1beta1 "github.com/skpr/operator/pkg/apis/aws/v1beta1"
)

// Helper function to get the active certificate request.
func getActiveRequestStatus(desired awsv1beta1.CertificateRequest, list []awsv1beta1.CertificateRequest) awsv1beta1.CertificateRequestReference {
	// Check the status of the desired certificate.
	if desired.Status.State == acm.CertificateStatusIssued {
		return requestToReference(desired)
	}

	// Fallback to the most recent "ISSUED" certificate.
	for _, item := range list {
		if item.Status.State == acm.CertificateStatusIssued {
			return requestToReference(item)
		}
	}

	// We didn't find a certificate which was "ISSUED".
	return awsv1beta1.CertificateRequestReference{}
}

// Helper function to convert a Request into a Reference.
func requestToReference(request awsv1beta1.CertificateRequest) awsv1beta1.CertificateRequestReference {
	return awsv1beta1.CertificateRequestReference{
		Name:    request.ObjectMeta.Name,
		Details: request.Status,
	}
}

// Helper function to get a hash from the certificate request spec.
func getHash(spec awsv1beta1.CertificateRequestSpec) (string, error) {
	// We have to sort the alternative names otherwise a new hash will
	// get generated when an order was changed.
	sort.Strings(spec.AlternateNames)

	return hashify.SHA1String(spec)
}

// Helper function to sort requests.
func sortRequests(list *awsv1beta1.CertificateRequestList) {
	sort.Slice(list.Items, func(i, j int) bool {
		return list.Items[i].ObjectMeta.CreationTimestamp.After(list.Items[j].ObjectMeta.CreationTimestamp.Time)
	})
}
