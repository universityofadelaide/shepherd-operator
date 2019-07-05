package certificate

import (
	"github.com/aws/aws-sdk-go/service/acm"
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

// Returns a list of CertificateRequests which have a common owner.
func filterByOwner(owner *awsv1beta1.Certificate, requests []awsv1beta1.CertificateRequest) []awsv1beta1.CertificateRequest {
	var list []awsv1beta1.CertificateRequest

	for _, request := range requests {
		for _, reference := range request.ObjectMeta.OwnerReferences {
			if reference.UID == owner.ObjectMeta.UID {
				list = append(list, request)
			}
		}
	}

	return list
}
