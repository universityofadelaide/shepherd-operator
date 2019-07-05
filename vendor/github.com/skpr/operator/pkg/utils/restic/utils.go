package restic

import "fmt"

// Helper function to format tags.
func formatTags(tags []string) string {
	var line string

	for _, tag := range tags {
		line = fmt.Sprintf("%s --tag=%s", line, tag)
	}

	return line
}
