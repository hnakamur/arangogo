package arangogo

import "fmt"

type DB struct {
	name string
	conn *Connection
}

func (c *Connection) DB(name string) *DB {
	return &DB{name: name, conn: c}
}

func (c *Connection) SystemDB() *DB {
	return c.DB(defaultDatabaseName)
}

func (c *Connection) CreateDatabase(name string, users []interface{}) error {
	payload := struct {
		Name  string        `json:"name"`
		Users []interface{} `json:"users"`
	}{
		Name: name,
	}
	if len(users) > 0 {
		payload.Users = users
	}
	body := new(struct {
		Error bool `json:"error"`
		Code  int  `json:"code"`
	})
	resp, err := c.send("POST", "/_api/database", payload, body)
	if err != nil {
		return fmt.Errorf("failed to create database: %v", err)
	}
	if body.Error {
		return fmt.Errorf("error in create database response:%s", string(resp.body))
	}
	return nil
}

func (c *Connection) DropDatabase(name string) error {
	body := new(struct {
		Error bool `json:"error"`
	})
	resp, err := c.send("DELETE", "/_api/database/"+name, nil, body)
	if err != nil {
		return fmt.Errorf("failed to delete database: %v", err)
	}
	if body.Error {
		return fmt.Errorf("error in delete database response:%s", string(resp.body))
	}
	return nil
}

func (c *Connection) ListDatabases() ([]string, error) {
	body := new(struct {
		Result []string `json:"result"`
		Error  bool     `json:"error"`
	})
	resp, err := c.send("GET", "/_api/database", nil, body)
	if err != nil {
		return nil, fmt.Errorf("failed to get database list: %v", err)
	}
	if body.Error {
		return nil, fmt.Errorf("error in database list response:%s", string(resp.body))
	}
	return body.Result, nil
}

func (c *Connection) ListUserDatabases() ([]string, error) {
	body := new(struct {
		Result []string `json:"result"`
		Error  bool     `json:"error"`
	})
	resp, err := c.send("GET", "/_api/database/user", nil, body)
	if err != nil {
		return nil, fmt.Errorf("failed to get user database list: %v", err)
	}
	if body.Error {
		return nil, fmt.Errorf("error in user database list response:%s", string(resp.body))
	}
	return body.Result, nil
}
