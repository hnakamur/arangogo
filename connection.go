package arangogo

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
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

const SystemDatabaseName = "_system"

func NewConnection(config *Config) (*Connection, error) {
	c := &Connection{
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
	return c, nil
}

type HTTPError struct {
	error
	StatusCode int
}

func (c *Connection) send(method, path string, header http.Header, payload, respBody interface{}) (*http.Response, error) {
	var reader io.Reader
	var payloadBytes []byte
	if payload != nil {
		var err error
		payloadBytes, err = json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("failed to encode request payload: %v", err)
		}
		reader = bytes.NewBuffer(payloadBytes)
	}
	url := c.url + path
	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	if c.header != nil {
		req.Header = c.header
	}
	if header != nil {
		if req.Header == nil {
			req.Header = header
		} else {
			for k, vv := range header {
				for _, v := range vv {
					req.Header.Add(k, v)
				}
			}
		}
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
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(b))
	if true {
		var reqH []byte
		for k, vv := range req.Header {
			for _, v := range vv {
				reqH = append(reqH, ", "+k+"="+v...)
			}
		}
		var payloadStr string
		if len(payloadBytes) > 0 {
			payloadStr = string(payloadBytes)
		}

		var respH []byte
		for k, vv := range resp.Header {
			for _, v := range vv {
				respH = append(respH, ", "+k+"="+v...)
			}
		}
		var bodyStr string
		if len(b) > 0 {
			bodyStr = string(b)
		}
		log.Printf("Connection send. method=%s, url=%s%s, payload=%s, status=%s%s, resBody=%s", method, url, string(reqH), payloadStr, resp.Status, respH, bodyStr)
	}
	if len(b) > 0 {
		errBody := new(struct {
			Error        bool   `json:"error"`
			ErrorNum     int    `json:"errorNum"`
			ErrorMessage string `json:"errorMessage"`
		})
		err2 := json.Unmarshal(b, errBody)
		if err2 == nil && errBody.Error {
			msg := fmt.Sprintf("error from ArangoDB. status=%d", resp.StatusCode)
			if errBody.ErrorNum != 0 {
				msg += fmt.Sprintf(", errorNum=%d", errBody.ErrorNum)
			}
			if errBody.ErrorMessage != "" {
				msg += fmt.Sprintf(", errorMessage=%s", errBody.ErrorMessage)
			}
			msg += fmt.Sprintf(", method=%s, url=%s", method, url)
			return nil, errors.New(msg)
		}

		if respBody != nil {
			err = json.Unmarshal(b, respBody)
			if err != nil {
				return nil, fmt.Errorf("failed to decode response body: %v", err)
			}
		}
	}
	s := resp.StatusCode
	if s >= http.StatusBadRequest {
		var bodyStr string
		if len(b) > 0 {
			bodyStr = string(b)
		}
		return nil, HTTPError{
			error:      fmt.Errorf("http status error:%s, body:%s", resp.Status, bodyStr),
			StatusCode: s,
		}
	}

	return resp, nil
}

type pathConfig struct {
	dbName      string
	pathFormat  string
	pathParams  []interface{}
	queryParams url.Values
}

func buildPath(c pathConfig) string {
	var path string
	if c.dbName != SystemDatabaseName && c.dbName != "" {
		path = "/_db/" + c.dbName
	}
	path += fmt.Sprintf(c.pathFormat, c.pathParams...)
	if len(c.queryParams) > 0 {
		path += "?" + c.queryParams.Encode()
	}
	return path
}
