package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommand(t *testing.T) {
	params := CommandParams{
		Service:   "aws s3",
		Operation: "cp",
		Args: []string{
			"foo.txt",
			"s3://bar/foo.txt",
		},
	}

	want := []string{
		"aws s3",
		"cp",
		"foo.txt",
		"s3://bar/foo.txt",
	}

	assert.Equal(t, want, Command(params))
}

func TestCommandWithEndpoint(t *testing.T) {
	params := CommandParams{
		Endpoint:  "http://localhost:9000",
		Service:   "aws s3",
		Operation: "cp",
		Args: []string{
			"foo.txt",
			"s3://bar/foo.txt",
		},
	}

	want := []string{
		"aws s3",
		"--endpoint-url", "http://localhost:9000",
		"cp",
		"foo.txt",
		"s3://bar/foo.txt",
	}

	assert.Equal(t, want, Command(params))
}
