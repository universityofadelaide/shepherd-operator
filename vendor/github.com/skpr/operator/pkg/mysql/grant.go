package mysql

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// GrantClient for interacting with MySQL.
type GrantClient struct {
	client *sql.DB
}

// Grant returns the grant client.
func (m *Client) Grant() *GrantClient {
	return &GrantClient{m.client}
}

// List returns a list of grants.
func (d GrantClient) List(username string) ([]string, error) {
	var grants []string

	rows, err := d.client.Query(fmt.Sprintf("SELECT Db FROM mysql.db WHERE User = '%s'", username))
	if err != nil {
		return grants, errors.Wrap(err, "query failed")
	}

	columns, err := rows.Columns()
	if err != nil {
		return grants, errors.Wrap(err, "failed to get column list")
	}

	// Make a slice for the values
	values := make([]sql.RawBytes, len(columns))

	// rows.Scan wants '[]interface{}' as an argument, so we must copy the
	// references into such a slice
	// See http://code.google.com/p/go-wiki/wiki/InterfaceSlice for details
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	// Fetch rows
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return grants, errors.Wrap(err, "failed to scan row")
		}

		for _, col := range values {
			// Here we can check if the value is nil (NULL value)
			if col == nil {
				continue
			}

			grants = append(grants, string(col))
		}
	}
	if err = rows.Err(); err != nil {
		return grants, errors.Wrap(err, "failed to get rows")
	}

	return grants, nil
}

// Exists checks if the grant exists.
func (d GrantClient) Exists(username, database string) (bool, error) {
	list, err := d.List(username)
	if err != nil {
		return false, errors.Wrap(err, "failed to list grants")
	}

	for _, item := range list {
		if item == database {
			return true, nil
		}
	}

	return false, nil
}

// Create a database and error if it exists.
func (d GrantClient) Create(username, database string, privileges []string) error {
	priv := strings.Join(privileges, ", ")

	_, err := d.client.Exec(fmt.Sprintf("GRANT %s ON %s.* To '%s'@'%%'", priv, database, username))
	if err != nil {
		return errors.Wrap(err, "failed to create grant")
	}

	return nil
}

// Revoke grant if it exists.
func (d GrantClient) Revoke(username, database string) error {
	_, err := d.client.Exec(fmt.Sprintf("REVOKE ALL ON %s.* FROM '%s'@'%%'", database, username))
	if err != nil {
		return errors.Wrap(err, "failed to revoke grant")
	}

	return nil
}
