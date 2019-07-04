package restic

import (
	"fmt"
	"regexp"
)

// Helper function to format tags.
func formatTags(tags []string) string {
	var line string

	for _, tag := range tags {
		line = fmt.Sprintf("%s --tag=%s", line, tag)
	}

	return line
}

// Helper function for parsing the snapshot id from restic backup output.
func parseSnapshotID(input string) string {
	// Restic IDs are SHA-256 hashes and the output contains the 8 character short version.
	var r = regexp.MustCompile(`snapshot\s([A-Fa-f0-9]{8})\ssaved`)
	match := r.FindStringSubmatch(input)
	if len(match) <= 1 {
		return ""
	}
	return match[1]
}
