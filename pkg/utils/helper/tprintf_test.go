package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Scenario struct {
	Format   string
	Args     map[string]interface{}
	Expected string
}

func TestTprintf(t *testing.T) {
	scenarios := []Scenario{
		{
			Format: "mysqldump --user={{.user}} --password={{.password}} --host={{.host}} --port={{.port}} --database={{.database}}",
			Args: map[string]interface{}{
				"user":     "bob",
				"password": "abc123",
				"host":     "mysql.svc",
				"port":     "3306",
				"database": "bob_prod",
			},
			Expected: "mysqldump --user=bob --password=abc123 --host=mysql.svc --port=3306 --database=bob_prod",
		},
		{
			Format: "kubectl -n {{.namespace}} get pod -o {{.format}}",
			Args: map[string]interface{}{
				"namespace": "default",
				"format":    "json",
			},
			Expected: "kubectl -n default get pod -o json",
		},
	}

	for _, scenario := range scenarios {
		actual, err := Tprintf(scenario.Format, scenario.Args)
		assert.Nil(t, err, "no errors encountered")
		assert.Equal(t, scenario.Expected, actual, "rendered string with named args correctly")
	}
}
