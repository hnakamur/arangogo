package arangogo

import (
	"fmt"
	"net/http"
)

type CreateDatabaseConfigUser struct {
	Username string `json:"username,omitempty"`
	Passwd   string `json:"passwd,omitempty"`
	Active   *bool  `json:"active,omitempty"`
}

type CreateDatabaseConfig struct {
	Username string                     `json:"username,omitempty"`
	Name     string                     `json:"name"`
	Extra    interface{}                `json:"extra,omitempty"`
	Passwd   string                     `json:"passwd,omitempty"`
	Active   *bool                      `json:"active,omitempty"`
	Users    []CreateDatabaseConfigUser `json:"users,omitempty"`
}

func (c *Connection) CreateDatabase(config CreateDatabaseConfig) error {
	_, err := c.send(http.MethodPost, "/_api/database", nil, config, nil)
	if err != nil {
		return fmt.Errorf("failed to create database: %v", err)
	}
	return nil
}

func (c *Connection) DropDatabase(name string) error {
	_, err := c.send(http.MethodDelete, "/_api/database/"+name, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to delete database: %v", err)
	}
	return nil
}

func (c *Connection) ListDatabases() ([]string, error) {
	var body struct {
		Result []string `json:"result"`
	}
	_, err := c.send(http.MethodGet, "/_api/database", nil, nil, &body)
	if err != nil {
		return nil, fmt.Errorf("failed to get database list: %v", err)
	}
	return body.Result, nil
}

func (c *Connection) ListUserDatabases() ([]string, error) {
	var body struct {
		Result []string `json:"result"`
	}
	_, err := c.send(http.MethodGet, "/_api/database/user", nil, nil, &body)
	if err != nil {
		return nil, fmt.Errorf("failed to get user database list: %v", err)
	}
	return body.Result, nil
}
