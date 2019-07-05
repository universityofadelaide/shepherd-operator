package cloudfront

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudfront"

	awsv1beta1 "github.com/skpr/operator/pkg/apis/aws/v1beta1"
)

const (
	// OriginID which maps behaviors to origins.
	OriginID = "skpr"
)

// Helper function to generate a CloudFront distribution.
// This object is really large due to required fields when updating the Distribution eg. FieldLevelEncryptionId
func generateDistribution(prefix string, instance *awsv1beta1.CloudFront) (*cloudfront.DistributionConfig, error) {
	reference := fmt.Sprintf("%s-%s-%s", prefix, instance.ObjectMeta.Namespace, instance.ObjectMeta.Name)

	behavior := &cloudfront.DefaultCacheBehavior{
		Compress:               aws.Bool(true),
		MinTTL:                 aws.Int64(0),
		MaxTTL:                 aws.Int64(86400),
		DefaultTTL:             aws.Int64(86400),
		TargetOriginId:         aws.String(OriginID),
		SmoothStreaming:        aws.Bool(true),
		FieldLevelEncryptionId: aws.String(""),
		ForwardedValues: &cloudfront.ForwardedValues{
			QueryString: aws.Bool(true),
			Headers: &cloudfront.Headers{
				Items:    aws.StringSlice(instance.Spec.Behavior.Whitelist.Headers),
				Quantity: aws.Int64(int64(len(instance.Spec.Behavior.Whitelist.Headers))),
			},
			QueryStringCacheKeys: &cloudfront.QueryStringCacheKeys{
				Quantity: aws.Int64(0),
			},
		},
		LambdaFunctionAssociations: &cloudfront.LambdaFunctionAssociations{
			Quantity: aws.Int64(0),
		},
		AllowedMethods: &cloudfront.AllowedMethods{
			Items: aws.StringSlice([]string{
				cloudfront.MethodHead,
				cloudfront.MethodDelete,
				cloudfront.MethodPost,
				cloudfront.MethodGet,
				cloudfront.MethodOptions,
				cloudfront.MethodPut,
				cloudfront.MethodPatch,
			}),
			Quantity: aws.Int64(7),
			CachedMethods: &cloudfront.CachedMethods{
				Items: aws.StringSlice([]string{
					cloudfront.MethodHead,
					cloudfront.MethodGet,
				}),
				Quantity: aws.Int64(2),
			},
		},
		TrustedSigners: &cloudfront.TrustedSigners{
			Enabled:  aws.Bool(false),
			Quantity: aws.Int64(0),
		},
		// This is a default and is overridden if an ACM certificate is specified.
		ViewerProtocolPolicy: aws.String(cloudfront.ViewerProtocolPolicyAllowAll),
	}

	if len(instance.Spec.Behavior.Whitelist.Cookies) > 0 {
		behavior.ForwardedValues.Cookies = &cloudfront.CookiePreference{
			Forward: aws.String(cloudfront.ItemSelectionWhitelist),
			WhitelistedNames: &cloudfront.CookieNames{
				Items:    aws.StringSlice(instance.Spec.Behavior.Whitelist.Cookies),
				Quantity: aws.Int64(int64(len(instance.Spec.Behavior.Whitelist.Cookies))),
			},
		}
	}

	origin := &cloudfront.Origin{
		Id:         aws.String(OriginID),
		DomainName: aws.String(instance.Spec.Origin.Endpoint),
		OriginPath: aws.String(""),
		CustomOriginConfig: &cloudfront.CustomOriginConfig{
			HTTPSPort:              aws.Int64(443),
			HTTPPort:               aws.Int64(80),
			OriginProtocolPolicy:   aws.String(instance.Spec.Origin.Policy),
			OriginReadTimeout:      aws.Int64(instance.Spec.Origin.Timeout),
			OriginKeepaliveTimeout: aws.Int64(5),
			OriginSslProtocols: &cloudfront.OriginSslProtocols{
				Items: aws.StringSlice([]string{
					cloudfront.SslProtocolTlsv1,
					cloudfront.SslProtocolTlsv11,
					cloudfront.SslProtocolTlsv12,
				}),
				Quantity: aws.Int64(3),
			},
		},
		CustomHeaders: &cloudfront.CustomHeaders{
			Quantity: aws.Int64(0),
		},
	}

	distribution := &cloudfront.DistributionConfig{
		Enabled:         aws.Bool(true),
		CallerReference: aws.String(reference),
		Comment:         aws.String("Automatically provisioned by github.com/skpr/operator"),
		IsIPV6Enabled:   aws.Bool(true),
		PriceClass:      aws.String(cloudfront.PriceClassPriceClassAll),
		Restrictions: &cloudfront.Restrictions{
			GeoRestriction: &cloudfront.GeoRestriction{
				Quantity:        aws.Int64(0),
				RestrictionType: aws.String(cloudfront.GeoRestrictionTypeNone),
			},
		},
		HttpVersion: aws.String(cloudfront.HttpVersionHttp2),
		Aliases: &cloudfront.Aliases{
			Quantity: aws.Int64(int64(len(instance.Spec.Aliases))),
		},
		DefaultRootObject:    aws.String(""),
		DefaultCacheBehavior: behavior,
		CacheBehaviors: &cloudfront.CacheBehaviors{
			Quantity: aws.Int64(0),
		},
		CustomErrorResponses: &cloudfront.CustomErrorResponses{
			Items: []*cloudfront.CustomErrorResponse{
				{
					ErrorCode:          aws.Int64(500),
					ErrorCachingMinTTL: aws.Int64(0),
					ResponseCode:       aws.String(""),
					ResponsePagePath:   aws.String(""),
				},
				{
					ErrorCode:          aws.Int64(502),
					ErrorCachingMinTTL: aws.Int64(0),
					ResponseCode:       aws.String(""),
					ResponsePagePath:   aws.String(""),
				},
				{
					ErrorCode:          aws.Int64(503),
					ErrorCachingMinTTL: aws.Int64(0),
					ResponseCode:       aws.String(""),
					ResponsePagePath:   aws.String(""),
				},
				{
					ErrorCode:          aws.Int64(504),
					ErrorCachingMinTTL: aws.Int64(0),
					ResponseCode:       aws.String(""),
					ResponsePagePath:   aws.String(""),
				},
			},
			Quantity: aws.Int64(4),
		},
		Origins: &cloudfront.Origins{
			Items: []*cloudfront.Origin{
				origin,
			},
			Quantity: aws.Int64(1),
		},
		OriginGroups: &cloudfront.OriginGroups{
			Quantity: aws.Int64(0),
		},
		Logging: &cloudfront.LoggingConfig{
			Enabled:        aws.Bool(false),
			Bucket:         aws.String(""), // @todo, Automatically provision bucket.
			Prefix:         aws.String(""), // @todo, Automatically provision bucket.
			IncludeCookies: aws.Bool(false),
		},
		WebACLId: aws.String(instance.Spec.Firewall.ARN),
		ViewerCertificate: &cloudfront.ViewerCertificate{
			CertificateSource:            aws.String(cloudfront.CertificateSourceCloudfront),
			CloudFrontDefaultCertificate: aws.Bool(true),
			MinimumProtocolVersion:       aws.String(cloudfront.MinimumProtocolVersionTlsv1),
		},
	}

	if len(instance.Spec.Aliases) > 0 {
		distribution.Aliases.Items = aws.StringSlice(instance.Spec.Aliases)
	}

	if instance.Spec.Certificate.ARN != "" {
		distribution.ViewerCertificate = &cloudfront.ViewerCertificate{
			Certificate:            aws.String(instance.Spec.Certificate.ARN),
			ACMCertificateArn:      aws.String(instance.Spec.Certificate.ARN),
			CertificateSource:      aws.String(cloudfront.CertificateSourceAcm),
			MinimumProtocolVersion: aws.String(cloudfront.MinimumProtocolVersionTlsv1),
			SSLSupportMethod:       aws.String(cloudfront.SSLSupportMethodSniOnly),
		}

		// Enforce HTTPS now that we have a certificate.
		distribution.DefaultCacheBehavior.ViewerProtocolPolicy = aws.String(cloudfront.ViewerProtocolPolicyRedirectToHttps)
	}

	return distribution, nil
}
