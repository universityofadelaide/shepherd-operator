package cloudfront

import (
	edgev1beta1 "github.com/skpr/operator/pkg/apis/edge/v1beta1"
	"github.com/skpr/operator/pkg/utils/slice"
)

// Helper function to get a list of domains from a list of routes.
func getDomainsFromRoutes(routes []edgev1beta1.IngressSpecRoute) []string {
	var domains []string

	for _, route := range routes {
		if slice.Contains(domains, route.Domain) {
			continue
		}

		domains = append(domains, route.Domain)
	}

	return domains
}
