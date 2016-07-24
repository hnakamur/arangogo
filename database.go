package arangogo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Config struct {
	URL           string
	DatabaseName  string
	ArangoVersion int
	Username      string
	Password      string
	Header        http.Header
}

type Connection struct {
	client        *http.Client
	url           string
	name          string
	arangoVersion int
	username      string
	password      string
	header        http.Header
}

const (
	defaultURL          = "http://localhost:8529"
	defaultDatabaseName = "_system"
	defaultArangoVesion = 30000
)

func NewConnection(config *Config) (*Connection, error) {
	d := &Connection{
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
			d.url = config.URL
		}
		if config.DatabaseName != "" {
			d.name = config.DatabaseName
		}
		if config.ArangoVersion != 0 {
			d.arangoVersion = config.ArangoVersion
		}
		if config.Username != "" {
			d.username = config.Username
		}
		if config.Password != "" {
			d.password = config.Password
		}
		if config.Header != nil {
			d.header = config.Header
		}
	}
	return d, nil
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

type HTTPError struct {
	error
	StatusCode int
}

func (c *Connection) send(method, path string, payload, respBody interface{}) (*response, error) {
	var reader io.Reader
	if payload != nil {
		b, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("failed to encode request payload: %v", err)
		}
		reader = bytes.NewBuffer(b)
	}
	url := c.url + path
	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	if c.header != nil {
		req.Header = c.header
	}
	if c.username != "" && c.password != "" {
		req.SetBasicAuth(c.username, c.password)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}
	if b != nil && respBody != nil {
		err = json.Unmarshal(b, respBody)
		if err != nil {
			return nil, fmt.Errorf("failed to decode response body: %v", err)
		}
	}
	s := resp.StatusCode
	if s >= http.StatusBadRequest {
		var bodyStr string
		if b != nil {
			bodyStr = string(b)
		}
		return nil, HTTPError{
			error:      fmt.Errorf("http status error:%s, body:%s", resp.Status, bodyStr),
			StatusCode: s,
		}
	}

	return &response{
		rawResponse: resp,
		body:        b,
	}, nil
}

type response struct {
	rawResponse *http.Response
	body        []byte
}

func (r *response) Status() string {
	return r.rawResponse.Status
}

func (r *response) StatusCode() int {
	return r.rawResponse.StatusCode
}
