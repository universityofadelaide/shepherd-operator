package metadata

// HasLabelWithValue checks if a label group has a key and value set.
func HasLabelWithValue(labels map[string]string, key, value string) bool {
	if _, ok := labels[key]; !ok {
		return false
	}

	if labels[key] != value {
		return false
	}

	return true
}
