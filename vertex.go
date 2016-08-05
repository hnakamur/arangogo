package arangogo

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type CreateVertexConfig struct {
	WaitForSync *bool
}

func (c *CreateVertexConfig) urlValues() url.Values {
	if c == nil {
		return nil
	}
	if c.WaitForSync != nil {
		return url.Values{
			"waitForSync": []string{
				strconv.FormatBool(*c.WaitForSync),
			},
		}
	}
	return nil
}

func (c *Connection) CreateVertex(dbName, graphName, collName string, data interface{}, config *CreateVertexConfig) (idKeyRev DocIDKeyRev, rc int, err error) {
	path := buildPath(pathConfig{
		dbName:      dbName,
		pathFormat:  "/_api/gharial/%s/vertex/%s",
		pathParams:  []interface{}{graphName, collName},
		queryParams: config.urlValues(),
	})

	var body struct {
		Vertex DocIDKeyRev `json:"vertex"`
		Code   int         `json:"code"`
	}
	_, err = c.send("POST", path, nil, data, &body)
	if err != nil {
		return body.Vertex, 0, fmt.Errorf("failed to create vertex: %v", err)
	}
	return body.Vertex, body.Code, nil
}

type GetVertexConfig struct {
	IfMatch string
}

func (c *GetVertexConfig) header() http.Header {
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

func (c *Connection) GetVertex(dbName, graphName, collName, vertexKey string, config *GetVertexConfig) (vertex interface{}, rc int, err error) {
	path := buildPath(pathConfig{
		dbName:     dbName,
		pathFormat: "/_api/gharial/%s/vertex/%s/%s",
		pathParams: []interface{}{graphName, collName, vertexKey},
	})

	var body struct {
		Vertex interface{} `json:"vertex"`
		Code   int         `json:"code"`
	}
	_, err = c.send("GET", path, config.header(), nil, &body)
	if err != nil {
		return body.Vertex, 0, fmt.Errorf("failed to create vertex: %v", err)
	}
	return body.Vertex, body.Code, nil
}
