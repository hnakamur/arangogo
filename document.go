package arangogo

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type DocIDKeyRev struct {
	ID  string `json:"_id"`
	Key string `json:"_key"`
	Rev string `json:"_rev"`
}

type CreateDocumentConfig struct {
	WaitForSync *bool
	ReturnNew   *bool
}

type Document struct {
	ID  string `json:"_id"`
	Key string `json:"_key"`
	Rev string `json:"_rev"`
}

func (c *Connection) CreateDocument(dbName, collName string, data interface{}, config *CreateDocumentConfig) (*Document, error) {
	body := new(Document)
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
	_, err := c.send("POST", u, nil, data, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create document: %v", err)
	}
	return body, nil
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
