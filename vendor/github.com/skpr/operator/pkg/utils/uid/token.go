package uid

import (
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/types"
)

// GetToken converts a Kubernetes UID into a 32 character which can be used as a token.
//   eg. AWS Certificate Requests require a 32 character idempotency token.
func GetToken(uid types.UID) (string, error) {
	token := strings.ReplaceAll(string(uid), "-", "")

	if len(token) > 32 {
		return "", fmt.Errorf("token is greater than 32 characters: %s", token)
	}

	return token, nil
}
