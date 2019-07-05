package hash

import (
	"crypto/md5"
	"fmt"
)

// String hash of the input.
func String(input string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(input)))
}
