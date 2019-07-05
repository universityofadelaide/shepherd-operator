package ses

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

// PasswordFromSecretKey was taken from the Terraform project.
// https://github.com/terraform-providers/terraform-provider-aws/blob/904a7762c5c6781c5429372aa79e8e6d84feed61/aws/resource_aws_iam_access_key.go#L199
func PasswordFromSecretKey(key string) (string, error) {
	if key == "" {
		return "", nil
	}

	version := byte(0x02)
	message := []byte("SendRawEmail")
	hmacKey := []byte(key)

	h := hmac.New(sha256.New, hmacKey)
	if _, err := h.Write(message); err != nil {
		return "", err
	}

	rawSig := h.Sum(nil)
	versionedSig := make([]byte, 0, len(rawSig)+1)
	versionedSig = append(versionedSig, version)
	versionedSig = append(versionedSig, rawSig...)

	return base64.StdEncoding.EncodeToString(versionedSig), nil
}
