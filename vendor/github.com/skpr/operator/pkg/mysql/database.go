package mysql

import (
	"database/sql"
	"fmt"

	"github.com/pkg/errors"
)

// DatabaseClient for interacting with
type DatabaseClient struct {
	client *sql.DB
}

// Database returns the Database client.
func (m *Client) Database() *DatabaseClient {
	return &DatabaseClient{m.client}
}

// List returns a list of databases.
func (d DatabaseClient) List() ([]string, error) {
	var databases []string

	rows, err := d.client.Query("SHOW DATABASES")
	if err != nil {
		return databases, errors.Wrap(err, "query failed")
	}

	columns, err := rows.Columns()
	if err != nil {
		return databases, errors.Wrap(err, "failed to get column list")
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
			return databases, errors.Wrap(err, "failed to scan row")
		}

		for _, col := range values {
			// Here we can check if the value is nil (NULL value)
			if col == nil {
				continue
			}

			databases = append(databases, string(col))
		}
	}
	if err = rows.Err(); err != nil {
		return databases, errors.Wrap(err, "failed to get rows")
	}

	return databases, nil
}

// Exists returns if a database exists.
func (d DatabaseClient) Exists(name string) (bool, error) {
	list, err := d.List()
	if err != nil {
		return false, errors.Wrap(err, "failed to list databases")
	}

	for _, item := range list {
		if item == name {
			return true, nil
		}
	}

	return false, nil
}

// Create will create a database and fail if it already exists.
func (d DatabaseClient) Create(name string) error {
	_, err := d.client.Exec(fmt.Sprintf("CREATE DATABASE %s", name))
	if err != nil {
		return errors.Wrap(err, "failed to create database")
	}

	return nil
}

// Delete will delete a database if it exists.
func (d DatabaseClient) Delete(name string) error {
	_, err := d.client.Exec(fmt.Sprintf("DROP DATABASE %s", name))
	if err != nil {
		return errors.Wrap(err, "failed to drop database")
	}

	return nil
}
