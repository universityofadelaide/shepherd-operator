package slice

// Contains a string within a slice.
func Contains(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}

	return false
}

// Remove a string from a slice.
func Remove(slice []string, s string) []string {
	var result []string

	for _, item := range slice {
		if item == s {
			continue
		}

		result = append(result, item)
	}

	return result
}
