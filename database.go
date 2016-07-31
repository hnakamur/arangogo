package arangogo

import (
	"encoding/json"
	"fmt"
)

type CreateDatabaseConfig struct {
	Username string                   `json:"username,omitempty"`
	Name     string                   `json:"name"`
	Extra    json.RawMessage          `json:"extra,omitempty"`
	Passwd   string                   `json:"passwd,omitempty"`
	Active   *bool                    `json:"active,omitempty"`
	Users    []map[string]interface{} `json:"users,omitempty"`
}

func (c *Connection) CreateDatabase(config CreateDatabaseConfig) error {
	_, err := c.send("POST", "/_api/database", nil, config, nil)
	if err != nil {
		return fmt.Errorf("failed to create database: %v", err)
	}
	return nil
}

func (c *Connection) DropDatabase(name string) error {
	_, err := c.send("DELETE", "/_api/database/"+name, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to delete database: %v", err)
	}
	return nil
}

func (c *Connection) ListDatabases() ([]string, error) {
	body := new(struct {
		Result []string `json:"result"`
	})
	_, err := c.send("GET", "/_api/database", nil, nil, body)
	if err != nil {
		return nil, fmt.Errorf("failed to get database list: %v", err)
	}
	return body.Result, nil
}

func (c *Connection) ListUserDatabases() ([]string, error) {
	body := new(struct {
		Result []string `json:"result"`
	})
	_, err := c.send("GET", "/_api/database/user", nil, nil, body)
	if err != nil {
		return nil, fmt.Errorf("failed to get user database list: %v", err)
	}
	return body.Result, nil
}
