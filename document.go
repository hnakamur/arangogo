package arangogo

import (
	"fmt"
	"net/url"
	"strconv"
)

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
	body := new(struct {
		Document
		Error bool `json:"error"`
	})
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
	resp, err := c.send("POST", u, data, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create document: %v", err)
	}
	if body.Error {
		return nil, fmt.Errorf("error in create document response:%s", string(resp.body))
	}
	return &body.Document, nil
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
	_, err := c.send("POST", u, data, &body)
	if err != nil {
		return nil, fmt.Errorf("failed to create documents: %v", err)
	}
	//	if body.Error {
	//		return nil, fmt.Errorf("error in create documents response:%s", string(resp.body))
	//	}
	return body, nil
}
