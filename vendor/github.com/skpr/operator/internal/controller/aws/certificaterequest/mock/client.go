package mock

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/acm"
	"github.com/aws/aws-sdk-go/service/acm/acmiface"
)

// Client which mocks the CloudFront client.
type Client struct {
	acmiface.ACMAPI
	certificates map[string]*acm.CertificateDetail
}

// New mock CloudFront client.
func New() *Client {
	return &Client{
		certificates: make(map[string]*acm.CertificateDetail),
	}
}

// RequestCertificate mock.
func (m *Client) RequestCertificate(input *acm.RequestCertificateInput) (*acm.RequestCertificateOutput, error) {
	if val, ok := m.certificates[*input.IdempotencyToken]; ok {
		return &acm.RequestCertificateOutput{
			CertificateArn: val.CertificateArn,
		}, nil
	}

	detail := &acm.CertificateDetail{
		CertificateArn:          input.IdempotencyToken,
		DomainName:              input.DomainName,
		SubjectAlternativeNames: input.SubjectAlternativeNames,
		DomainValidationOptions: []*acm.DomainValidation{
			{
				ResourceRecord: &acm.ResourceRecord{
					Name:  aws.String("aaaaaaaaaaaaa"),
					Type:  aws.String("bbbbbbbbbbbbb"),
					Value: aws.String("ccccccccccccc"),
				},
				ValidationStatus: aws.String(acm.CertificateStatusIssued),
			},
			{
				ResourceRecord: &acm.ResourceRecord{
					Name:  aws.String("ddddddddddddd"),
					Type:  aws.String("eeeeeeeeeeeee"),
					Value: aws.String("fffffffffffff"),
				},
				ValidationStatus: aws.String(acm.CertificateStatusIssued),
			},
			{
				ResourceRecord: &acm.ResourceRecord{
					Name:  aws.String("ggggggggggggg"),
					Type:  aws.String("hhhhhhhhhhhhh"),
					Value: aws.String("iiiiiiiiiiiii"),
				},
				ValidationStatus: aws.String(acm.CertificateStatusIssued),
			},
		},
		Status: aws.String(acm.CertificateStatusIssued),
	}

	m.certificates[*detail.CertificateArn] = detail

	return &acm.RequestCertificateOutput{
		CertificateArn: detail.CertificateArn,
	}, nil
}

// DescribeCertificate mock.
func (m *Client) DescribeCertificate(input *acm.DescribeCertificateInput) (*acm.DescribeCertificateOutput, error) {
	if val, ok := m.certificates[*input.CertificateArn]; ok {
		return &acm.DescribeCertificateOutput{
			Certificate: val,
		}, nil
	}

	return nil, awserr.New(acm.ErrCodeResourceNotFoundException, "not found", errors.New("distribution not found"))
}

// DeleteCertificate mock.
func (m *Client) DeleteCertificate(input *acm.DeleteCertificateInput) (*acm.DeleteCertificateOutput, error) {
	resp := &acm.DeleteCertificateOutput{}

	if _, ok := m.certificates[*input.CertificateArn]; !ok {
		return resp, awserr.New(acm.ErrCodeResourceNotFoundException, "not found", errors.New("distribution not found"))
	}

	delete(m.certificates, *input.CertificateArn)

	return resp, nil
}
