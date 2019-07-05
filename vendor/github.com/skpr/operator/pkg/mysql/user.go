package mysql

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
)

// UserClient for interacting with MySQL.
type UserClient struct {
	client *sql.DB
}

// User returns the user client.
func (m *Client) User() *UserClient {
	return &UserClient{m.client}
}

// List returns a list of users.
func (d UserClient) List() ([]string, error) {
	var users []string

	rows, err := d.client.Query("SELECT User FROM mysql.user")
	if err != nil {
		return users, errors.Wrap(err, "query failed")
	}

	columns, err := rows.Columns()
	if err != nil {
		return users, errors.Wrap(err, "failed to get column list")
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
			return users, errors.Wrap(err, "failed to scan row")
		}

		for _, col := range values {
			// Here we can check if the value is nil (NULL value)
			if col == nil {
				continue
			}

			users = append(users, string(col))
		}
	}
	if err = rows.Err(); err != nil {
		return users, errors.Wrap(err, "failed to get rows")
	}

	return users, nil
}

// Exists check if the user has already been created.
func (d UserClient) Exists(name string) (bool, error) {
	list, err := d.List()
	if err != nil {
		return false, errors.Wrap(err, "failed to list Users")
	}

	for _, item := range list {
		if item == name {
			return true, nil
		}
	}

	return false, nil
}

// Create a user and error if it already exists.
func (d UserClient) Create(username, password string) error {
	_, err := d.client.Exec(fmt.Sprintf("CREATE USER '%s'@'%%' IDENTIFIED BY '%s'", username, password))
	if err != nil {
		return errors.Wrap(err, "failed to create User")
	}

	return nil
}

// Delete a user.
func (d UserClient) Delete(name string) error {
	_, err := d.client.Exec(fmt.Sprintf("DROP USER '%s'@'%%'", name))
	if err != nil {
		return errors.Wrap(err, "failed to delete grant")
	}

	return nil
}
