package shared

import (
	"strconv"

	"github.com/spaolacci/murmur3"
)

// Helper function to generate a deterministic hash.
func hash(incoming string) string {
	h64 := murmur3.New64()
	h64.Write([]byte(incoming))
	result := h64.Sum64()
	return strconv.FormatUint(result, 36)
}
