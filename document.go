package arangogo

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type Document struct {
	ID  string `json:"_id"`
	Key string `json:"_key"`
	Rev string `json:"_rev"`
}

type CreateDocumentConfig struct {
	WaitForSync *bool
	ReturnNew   *bool
}

func (c *CreateDocumentConfig) queryParams() url.Values {
	if c == nil {
		return nil
	}

	var params url.Values
	if c.WaitForSync != nil || c.ReturnNew != nil {
		params = make(url.Values)
	}
	if c.WaitForSync != nil {
		params.Set("waitForSync", strconv.FormatBool(*c.WaitForSync))
	}
	if c.ReturnNew != nil {
		params.Set("returnNew", strconv.FormatBool(*c.ReturnNew))
	}
	return params
}

func (c *Connection) CreateDocument(dbName, collName string, data interface{}, config *CreateDocumentConfig, docPtr interface{}) (doc Document, rc int, err error) {
	path := buildPath(pathConfig{
		dbName:      dbName,
		pathFormat:  "/_api/document/%s",
		pathParams:  []interface{}{collName},
		queryParams: config.queryParams(),
	})

	var body struct {
		ID  string      `json:"_id"`
		Key string      `json:"_key"`
		Rev string      `json:"_rev"`
		New interface{} `json:"new"`
	}
	if docPtr != nil {
		body.New = docPtr
	}
	resp, err := c.send(http.MethodPost, path, nil, data, &body)
	if err != nil {
		return doc, resp.rawResponse.StatusCode, fmt.Errorf("failed to create document: %v", err)
	}
	doc = Document{
		ID:  body.ID,
		Key: body.Key,
		Rev: body.Rev,
	}
	return doc, resp.rawResponse.StatusCode, nil
}

func (c *Connection) CreateDocuments(dbName, collName string, data interface{}, config *CreateDocumentConfig) ([]Document, error) {
	var body []Document
	u := dbPrefix(dbName) + "/_api/document/" + collName
	v := url.Values{}
	if config != nil {
		if config.WaitForSync != nil {
			v.Set("waitForSync", strconv.FormatBool(*config.WaitForSync))
		}
		if config.ReturnNew != nil {
			v.Set("returnNew", strconv.FormatBool(*config.ReturnNew))
		}
	}
	if len(v) > 0 {
		u = u + "?" + v.Encode()
	}
	_, err := c.send("POST", u, nil, data, &body)
	if err != nil {
		return nil, fmt.Errorf("failed to create documents: %v", err)
	}
	return body, nil
}

type DeleteDocumentConfig struct {
	WaitForSync *bool
	ReturnOld   *bool
	IfMatch     string
}

func (c *Connection) DeleteDocument(dbName, collName, key string, config *DeleteDocumentConfig) error {
	u := dbPrefix(dbName) + "/_api/document/" + collName + "/" + key
	v := url.Values{}
	var header http.Header
	if config != nil {
		if config.WaitForSync != nil {
			v.Set("waitForSync", strconv.FormatBool(*config.WaitForSync))
		}
		if config.ReturnOld != nil {
			v.Set("returnOld", strconv.FormatBool(*config.ReturnOld))
		}
		if config.IfMatch != "" {
			header = make(http.Header)
			header.Set("if-match", config.IfMatch)
		}
	}
	if len(v) > 0 {
		u = u + "?" + v.Encode()
	}
	_, err := c.send("DELETE", u, header, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to delete document: %v", err)
	}
	return nil
}
