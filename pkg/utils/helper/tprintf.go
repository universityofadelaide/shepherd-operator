package helper

import (
	"bytes"
	"text/template"
)

// Tprintf allows substitution of named arguments in the format string. Useful
// if the argument list for Sprintf() would be cumbersome to manage.
func Tprintf(format string, args map[string]interface{}) (string, error) {
	t := template.Must(template.New("").Parse(format))
	var buf bytes.Buffer
	err := t.Execute(&buf, args)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

// TprintfMustParse panics on errors encountered during render.
func TprintfMustParse(format string, args map[string]interface{}) string {
	rendered, err := Tprintf(format, args)
	if err != nil {
		panic(err)
	}
	return rendered
}
