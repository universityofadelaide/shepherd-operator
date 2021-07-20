package restic

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// Helper function to format tags.
func formatTags(tags []string) string {
	var line string

	for _, tag := range tags {
		line = fmt.Sprintf("%s --tag=%s", line, tag)
	}

	return strings.Trim(line, " ")
}

// ParseSnapshotID parses the restic snapshot id from a string.
func ParseSnapshotID(input string) string {
	// Restic IDs are SHA-256 hashes and the output contains the 8 character short version.
	var r = regexp.MustCompile(`snapshot\s([A-Fa-f0-9]{8})\ssaved`)
	match := r.FindStringSubmatch(input)
	if len(match) <= 1 {
		return ""
	}
	return match[1]
}

// RequeueAfterSeconds returns a reconcile.Result to requeue after seconds time.
func RequeueAfterSeconds(seconds int64) reconcile.Result {
	return reconcile.Result{
		Requeue:      true,
		RequeueAfter: time.Duration(seconds) * time.Second,
	}
}
