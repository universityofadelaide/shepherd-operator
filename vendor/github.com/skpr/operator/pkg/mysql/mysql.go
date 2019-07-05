package mysql

import (
	"database/sql"
	"fmt"
)

const (
	// PrivilegeAll allows a user to perform all operations on a database.
	PrivilegeAll = "ALL"
	// PrivilegeProcess allows a user to see process information.
	PrivilegeProcess = "PROCESS"
	// PrivilegeReplication allows a user to see replication information.
	PrivilegeReplication = "REPLICATION CLIENT"
	// PrivilegeSelect allows a user to perform select operations.
	PrivilegeSelect = "SELECT"
	// PrivilegeInsert allows for insert operations.
	PrivilegeInsert = "INSERT"
	// PrivilegeUpdate allows for update operations.
	PrivilegeUpdate = "UPDATE"
	// PrivilegeDelete allows for delete operations.
	PrivilegeDelete = "DELETE"
	// PrivilegeCreate allows for create operations.
	PrivilegeCreate = "CREATE"
	// PrivilegeDrop allows for drop  operations.
	PrivilegeDrop = "DROP"
	// PrivilegeIndex allows for index operations.
	PrivilegeIndex = "INDEX"
	// PrivilegeAlter allows for alter operations.
	PrivilegeAlter = "ALTER"
	// PrivilegeCrateTemporaryTables allows for temporary tables.
	PrivilegeCrateTemporaryTables = "CREATE TEMPORARY TABLES"
)

// Client for interacting with MySQL.
type Client struct {
	client *sql.DB
}

// New client for interfacting with MySQL.
func New(host, user, pass string, port int) (*Client, error) {
	conn := fmt.Sprintf("%s:%s@tcp(%s:%d)/", user, pass, host, port)

	client, err := sql.Open("mysql", conn)
	if err != nil {
		return nil, err
	}

	err = client.Ping()
	if err != nil {
		return nil, err
	}

	return &Client{client}, nil
}
