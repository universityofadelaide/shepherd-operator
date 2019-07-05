package policy

import (
	"encoding/json"
)

// Print the document.
func Print(document Document) (string, error) {
	raw, err := json.Marshal(document)
	if err != nil {
		return "", err
	}

	return string(raw), nil
}
