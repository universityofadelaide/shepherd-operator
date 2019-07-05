package prometheus

import "fmt"

// Path for scraping metrics.
func Path(path, token string) string {
	return fmt.Sprintf("%s?token=%s", path, token)
}
