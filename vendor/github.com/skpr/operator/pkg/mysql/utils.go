package mysql

import "strings"

// Format to ensure it is the correct format.
func Format(name string) string {
	name = strings.Replace(name, "-", "_", -1)
	return name
}
