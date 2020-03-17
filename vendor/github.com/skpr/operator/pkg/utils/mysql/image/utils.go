package image

// Helper function to format destinations.
func formatDestination(destinations []string) []string {
	var out []string

	for _, d := range destinations {
		out = append(out, "--destination", d)
	}

	return out
}
