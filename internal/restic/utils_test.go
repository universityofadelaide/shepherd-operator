//go:build unit
// +build unit

package restic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatTags(t *testing.T) {
	assert.Equal(t, "--tag=foo", formatTags([]string{"foo"}))
	assert.Equal(t, "--tag=foo --tag=bar", formatTags([]string{"foo", "bar"}))
}

func TestParseResticId(t *testing.T) {
	var input = `
[restic-backup] open repository
[restic-backup] created new cache in /root/.cache/restic
[restic-backup] lock repository
[restic-backup] load index files
[restic-backup] start scan on [.]
[restic-backup] start backup on [.]
[restic-backup] scan finished in 0.804s: 1 files, 1.308 KiB
[restic-backup] 
[restic-backup] Files:           1 new,     0 changed,     0 unmodified
[restic-backup] Dirs:            0 new,     0 changed,     0 unmodified
[restic-backup] Data Blobs:      1 new
[restic-backup] Tree Blobs:      1 new
[restic-backup] Added to the repo: 1.993 KiB
[restic-backup] 
[restic-backup] processed 1 files, 1.308 KiB in 0:00
[restic-backup] snapshot 12487f64 saved`
	assert.Equal(t, "12487f64", ParseSnapshotID(input))
	input = `
[restic-backup] open repository
[restic-backup] created new cache in /root/.cache/restic
[restic-backup] lock repository
[restic-backup] load index files
[restic-backup] start scan on [.]
[restic-backup] start backup on [.]
[restic-backup] scan finished in 0.804s: 1 files, 1.308 KiB
[restic-backup] 
[restic-backup] Files:           1 new,     0 changed,     0 unmodified
[restic-backup] Dirs:            0 new,     0 changed,     0 unmodified
[restic-backup] Data Blobs:      1 new
[restic-backup] Tree Blobs:      1 new
[restic-backup] Added to the repo: 1.993 KiB
[restic-backup] 
[restic-backup] processed 1 files, 1.308 KiB in 0:00`
	assert.Equal(t, "", ParseSnapshotID(input))
	assert.Equal(t, "", ParseSnapshotID("snapshot 12487f64ASD saved"))
	assert.Equal(t, "12345678", ParseSnapshotID("snapshot 12345678 saved"))
	assert.Equal(t, "1a3b5c7d", ParseSnapshotID("snapshot 1a3b5c7d saved"))
}
