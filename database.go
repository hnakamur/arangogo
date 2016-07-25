package arangogo

import (
	"fmt"
	"net/http"
	"net/url"
)

type Database struct {
	conn *connection
}

func NewDatabase(config *Config) (*Database, error) {
	c := &connection{
		client:        new(http.Client),
		url:           defaultURL,
		name:          defaultDatabaseName,
		arangoVersion: defaultArangoVesion,
	}
	if config != nil {
		if config.URL != "" {
			_, err := url.Parse(config.URL)
			if err != nil {
				return nil, err
			}
			c.url = config.URL
		}
		if config.DatabaseName != "" {
			c.name = config.DatabaseName
		}
		if config.ArangoVersion != 0 {
			c.arangoVersion = config.ArangoVersion
		}
		if config.Username != "" {
			c.username = config.Username
		}
		if config.Password != "" {
			c.password = config.Password
		}
		if config.Header != nil {
			c.header = config.Header
		}
	}
	return &Database{conn: c}, nil
}

func (d *Database) CreateDatabase(name string, users []interface{}) error {
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
	resp, err := d.conn.send("POST", "/_api/database", payload, body)
	if err != nil {
		return fmt.Errorf("failed to create database: %v", err)
	}
	if body.Error {
		return fmt.Errorf("error in create database response:%s", string(resp.body))
	}
	return nil
}

func (d *Database) DropDatabase(name string) error {
	body := new(struct {
		Error bool `json:"error"`
	})
	resp, err := d.conn.send("DELETE", "/_api/database/"+name, nil, body)
	if err != nil {
		return fmt.Errorf("failed to delete database: %v", err)
	}
	if body.Error {
		return fmt.Errorf("error in delete database response:%s", string(resp.body))
	}
	return nil
}

func (d *Database) ListDatabases() ([]string, error) {
	body := new(struct {
		Result []string `json:"result"`
		Error  bool     `json:"error"`
	})
	resp, err := d.conn.send("GET", "/_api/database", nil, body)
	if err != nil {
		return nil, fmt.Errorf("failed to get database list: %v", err)
	}
	if body.Error {
		return nil, fmt.Errorf("error in database list response:%s", string(resp.body))
	}
	return body.Result, nil
}

func (d *Database) ListUserDatabases() ([]string, error) {
	body := new(struct {
		Result []string `json:"result"`
		Error  bool     `json:"error"`
	})
	resp, err := d.conn.send("GET", "/_api/database/user", nil, body)
	if err != nil {
		return nil, fmt.Errorf("failed to get user database list: %v", err)
	}
	if body.Error {
		return nil, fmt.Errorf("error in user database list response:%s", string(resp.body))
	}
	return body.Result, nil
}

func (d *Database) ListCollections(excludeSystem bool) ([]Collection, error) {
	body := new(struct {
		Result []Collection `json:"result"`
		Error  bool         `json:"error"`
	})
	resp, err := d.conn.send("GET", "/_api/collection", nil, body)
	if err != nil {
		return nil, fmt.Errorf("failed to get collection list: %v", err)
	}
	if body.Error {
		return nil, fmt.Errorf("error in collection list response:%s", string(resp.body))
	}
	var collections []Collection
	if excludeSystem {
		for _, c := range body.Result {
			if !c.IsSystem {
				collections = append(collections, c)
			}
		}
	} else {
		collections = body.Result
	}
	return collections, nil
}
