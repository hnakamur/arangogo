package arangogo

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type ReadDocumentConfig struct {
	IfNoneMatch string
	IfMatch     string
}

func (c *ReadDocumentConfig) header() http.Header {
	if c == nil {
		return nil
	}
	var header http.Header
	if c.IfNoneMatch != "" || c.IfMatch != "" {
		header = make(http.Header)
	}
	if c.IfNoneMatch != "" {
		header.Set("if-none-match", c.IfNoneMatch)
	}
	if c.IfMatch != "" {
		header.Set("if-match", c.IfMatch)
	}
	return header
}

func (c *Connection) ReadDocument(dbName, documentHandle string, config *ReadDocumentConfig, documentPtr interface{}) (rc int, err error) {
	path := buildPath(pathConfig{
		dbName:     dbName,
		pathFormat: "/_api/document/%s",
		pathParams: []interface{}{documentHandle},
	})

	rc, _, err = c.send(http.MethodGet, path, config.header(), nil, &documentPtr)
	if err != nil {
		return rc, fmt.Errorf("failed to read document: %v", err)
	}
	return rc, nil
}

type ReadDocumentHeaderConfig struct {
	IfNoneMatch string
	IfMatch     string
}

func (c *ReadDocumentHeaderConfig) header() http.Header {
	if c == nil {
		return nil
	}
	var header http.Header
	if c.IfNoneMatch != "" || c.IfMatch != "" {
		header = make(http.Header)
	}
	if c.IfNoneMatch != "" {
		header.Set("if-none-match", c.IfNoneMatch)
	}
	if c.IfMatch != "" {
		header.Set("if-match", c.IfMatch)
	}
	return header
}

func (c *Connection) ReadDocumentHeader(dbName, documentHandle string, config *ReadDocumentHeaderConfig) (rev string, rc int, err error) {
	path := buildPath(pathConfig{
		dbName:     dbName,
		pathFormat: "/_api/document/%s",
		pathParams: []interface{}{documentHandle},
	})

	rc, resp, err := c.send(http.MethodHead, path, config.header(), nil, nil)
	rev = resp.Header.Get("ETag")
	if err != nil {
		return rev, rc, fmt.Errorf("failed to read document header: %v", err)
	}
	return rev, rc, nil
}

const (
	ListAllDocumentsTypeID   = "id"
	ListAllDocumentsTypeKey  = "key"
	ListAllDocumentsTypePath = "path"
)

type ListAllDocumentsConfig struct {
	Type       string `json:"type,omitempty"`
	Collection string `json:"collection"`
}

type ListAllDocumentsResultStats struct {
	WritesExecuted int     `json:"writesExecuted"`
	WritesIgnored  int     `json:"writesIgnored"`
	ScannedFull    int     `json:"scannedFull"`
	ScannedIndex   int     `json:"scannedIndex"`
	Filtered       int     `json:"filtered"`
	ExecutionTime  float64 `json:"executionTime"`
}

type ListAllDocumentsResult struct {
	Result  []string `json:"result"`
	HasMore bool     `json:"hasMore"`
	Cached  bool     `json:"cached"`
	Extra   struct {
		Stats ListAllDocumentsResultStats `json:"stats"`
	} `json:"extra"`
}

func (c *Connection) ListAllDocuments(dbName string, config ListAllDocumentsConfig) (r ListAllDocumentsResult, rc int, err error) {
	path := buildPath(pathConfig{
		dbName:     dbName,
		pathFormat: "/_api/simple/all-keys",
	})

	rc, _, err = c.send(http.MethodPut, path, nil, config, &r)
	if err != nil {
		err = fmt.Errorf("failed to list all documents: %v", err)
	}
	return
}

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
	rc, _, err = c.send(http.MethodPost, path, nil, data, &body)
	if err != nil {
		return doc, rc, fmt.Errorf("failed to create document: %v", err)
	}
	doc = Document{
		ID:  body.ID,
		Key: body.Key,
		Rev: body.Rev,
	}
	return doc, rc, nil
}

type CreateDocumentsConfig struct {
	WaitForSync *bool
}

func (c *CreateDocumentsConfig) queryParams() url.Values {
	if c == nil {
		return nil
	}

	var params url.Values
	if c.WaitForSync != nil {
		params = make(url.Values)
		params.Set("waitForSync", strconv.FormatBool(*c.WaitForSync))
	}
	return params
}

func (c *Connection) CreateDocuments(dbName, collName string, data interface{}, config *CreateDocumentsConfig) (docs []Document, rc int, err error) {
	path := buildPath(pathConfig{
		dbName:      dbName,
		pathFormat:  "/_api/document/%s",
		pathParams:  []interface{}{collName},
		queryParams: config.queryParams(),
	})

	rc, _, err = c.send(http.MethodPost, path, nil, data, &docs)
	if err != nil {
		return nil, rc, fmt.Errorf("failed to create documents: %v", err)
	}
	return docs, rc, nil
}

type ReplaceDocumentConfig struct {
	WaitForSync *bool
	IgnoreRevs  *bool
	ReturnOld   *bool
	ReturnNew   *bool
	IfMatch     string
}

func (c *ReplaceDocumentConfig) header() http.Header {
	if c == nil {
		return nil
	}
	var header http.Header
	if c.IfMatch != "" {
		header = make(http.Header)
		header.Set("if-match", c.IfMatch)
		return header
	}

	return nil
}

func (c *ReplaceDocumentConfig) queryParams() url.Values {
	if c == nil {
		return nil
	}

	var params url.Values
	if c.WaitForSync != nil || c.IgnoreRevs != nil || c.ReturnOld != nil || c.ReturnNew != nil {
		params = make(url.Values)
	}
	if c.WaitForSync != nil {
		params.Set("waitForSync", strconv.FormatBool(*c.WaitForSync))
	}
	if c.IgnoreRevs != nil {
		params.Set("ignoreRevs", strconv.FormatBool(*c.IgnoreRevs))
	}
	if c.ReturnOld != nil {
		params.Set("returnOld", strconv.FormatBool(*c.ReturnOld))
	}
	if c.ReturnNew != nil {
		params.Set("returnNew", strconv.FormatBool(*c.ReturnNew))
	}
	return params
}

type ReplaceDocumentResult struct {
	ID     string `json:"_id"`
	Key    string `json:"_key"`
	Rev    string `json:"_rev"`
	OldRev string `json:"_oldRev"`
}

func (c *Connection) ReplaceDocument(dbName, docHandle string, data interface{}, config *ReplaceDocumentConfig, oldDocPtr, newDocPtr interface{}) (r ReplaceDocumentResult, rc int, err error) {
	path := buildPath(pathConfig{
		dbName:      dbName,
		pathFormat:  "/_api/document/%s",
		pathParams:  []interface{}{docHandle},
		queryParams: config.queryParams(),
	})

	var body struct {
		ID     string      `json:"_id"`
		Key    string      `json:"_key"`
		Rev    string      `json:"_rev"`
		OldRev string      `json:"_oldRev"`
		Old    interface{} `json:"old"`
		New    interface{} `json:"new"`
	}
	if oldDocPtr != nil {
		body.Old = oldDocPtr
	}
	if newDocPtr != nil {
		body.New = newDocPtr
	}
	rc, _, err = c.send(http.MethodPut, path, config.header(), data, &body)
	if err != nil {
		return r, rc, fmt.Errorf("failed to replace document: %v", err)
	}
	r = ReplaceDocumentResult{
		ID:     body.ID,
		Key:    body.Key,
		Rev:    body.Rev,
		OldRev: body.OldRev,
	}
	return r, rc, nil
}

type RemoveDocumentConfig struct {
	WaitForSync *bool
	ReturnOld   *bool
	IfMatch     string
}

func (c *RemoveDocumentConfig) header() http.Header {
	if c == nil {
		return nil
	}
	var header http.Header
	if c.IfMatch != "" {
		header = make(http.Header)
		header.Set("if-match", c.IfMatch)
		return header
	}

	return nil
}

func (c *RemoveDocumentConfig) queryParams() url.Values {
	if c == nil {
		return nil
	}

	var params url.Values
	if c.WaitForSync != nil || c.ReturnOld != nil {
		params = make(url.Values)
	}
	if c.WaitForSync != nil {
		params.Set("waitForSync", strconv.FormatBool(*c.WaitForSync))
	}
	if c.ReturnOld != nil {
		params.Set("returnOld", strconv.FormatBool(*c.ReturnOld))
	}
	return params
}

func (c *Connection) RemoveDocument(dbName, collName, key string, config *RemoveDocumentConfig, docPtr interface{}) (doc Document, rc int, err error) {
	path := buildPath(pathConfig{
		dbName:      dbName,
		pathFormat:  "/_api/document/%s/%s",
		pathParams:  []interface{}{collName, key},
		queryParams: config.queryParams(),
	})

	var body struct {
		ID  string      `json:"_id"`
		Key string      `json:"_key"`
		Rev string      `json:"_rev"`
		Old interface{} `json:"old"`
	}
	if docPtr != nil {
		body.Old = docPtr
	}
	rc, _, err = c.send(http.MethodDelete, path, config.header(), nil, &body)
	if err != nil {
		return doc, rc, fmt.Errorf("failed to remove document: %v", err)
	}
	doc = Document{
		ID:  body.ID,
		Key: body.Key,
		Rev: body.Rev,
	}
	return doc, rc, nil
}
