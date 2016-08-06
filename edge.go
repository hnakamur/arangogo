package arangogo

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type CreateEdgeResult struct {
	ID  string `json:"_id"`
	Key string `json:"_key"`
	Rev string `json:"_rev"`
}

type CreateEdgeConfig struct {
	WaitForSync *bool
}

func (c *CreateEdgeConfig) urlValues() url.Values {
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

func (c *Connection) CreateEdge(dbName, graphName, collName string, data interface{}, config *CreateEdgeConfig) (r CreateEdgeResult, rc int, err error) {
	path := buildPath(pathConfig{
		dbName:      dbName,
		pathFormat:  "/_api/gharial/%s/edge/%s",
		pathParams:  []interface{}{graphName, collName},
		queryParams: config.urlValues(),
	})

	var body struct {
		Edge CreateEdgeResult `json:"edge"`
		Code int              `json:"code"`
	}
	_, err = c.send(http.MethodPost, path, nil, data, &body)
	if err != nil {
		return body.Edge, 0, fmt.Errorf("failed to create edge: %v", err)
	}
	return body.Edge, body.Code, nil
}

type GetEdgeConfig struct {
	IfMatch string
}

func (c *GetEdgeConfig) header() http.Header {
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

func (c *Connection) GetEdge(dbName, graphName, collName, edgeKey string, config *GetEdgeConfig, edgePtr interface{}) (rc int, err error) {
	path := buildPath(pathConfig{
		dbName:     dbName,
		pathFormat: "/_api/gharial/%s/edge/%s/%s",
		pathParams: []interface{}{graphName, collName, edgeKey},
	})

	var body struct {
		Edge interface{} `json:"edge"`
		Code int         `json:"code"`
	}
	if edgePtr != nil {
		body.Edge = edgePtr
	}
	_, err = c.send(http.MethodGet, path, config.header(), nil, &body)
	if err != nil {
		return 0, fmt.Errorf("failed to get edge: %v", err)
	}
	return body.Code, nil
}

type ModifyEdgeResult struct {
	ID     string `json:"_id"`
	Key    string `json:"_key"`
	Rev    string `json:"_rev"`
	OldRev string `json:"_oldRev"`
}

type ModifyEdgeConfig struct {
	WaitForSync *bool
	KeepNull    *bool
}

func (c *ModifyEdgeConfig) queryParams() url.Values {
	if c == nil {
		return nil
	}

	var params url.Values
	if c.WaitForSync != nil || c.KeepNull != nil {
		params = make(url.Values)
	}
	if c.WaitForSync != nil {
		params.Set("waitForSync", strconv.FormatBool(*c.WaitForSync))
	}
	if c.KeepNull != nil {
		params.Set("keepNull", strconv.FormatBool(*c.KeepNull))
	}
	return params
}

func (c *Connection) ModifyEdge(dbName, graphName, collName, edgeKey string, data interface{}, config *ModifyEdgeConfig) (edge ModifyEdgeResult, rc int, err error) {
	path := buildPath(pathConfig{
		dbName:      dbName,
		pathFormat:  "/_api/gharial/%s/edge/%s/%s",
		pathParams:  []interface{}{graphName, collName, edgeKey},
		queryParams: config.queryParams(),
	})

	var body struct {
		Edge ModifyEdgeResult `json:"edge"`
		Code int              `json:"code"`
	}
	_, err = c.send(http.MethodPatch, path, nil, data, &body)
	if err != nil {
		return body.Edge, 0, fmt.Errorf("failed to modify edge: %v", err)
	}
	return body.Edge, body.Code, nil
}

//type ReplaceVertexResult struct {
//	ID     string `json:"_id"`
//	Key    string `json:"_key"`
//	Rev    string `json:"_rev"`
//	OldRev string `json:"_oldRev"`
//}
//
//type ReplaceVertexConfig struct {
//	WaitForSync *bool
//	IfMatch     string
//}
//
//func (c *ReplaceVertexConfig) header() http.Header {
//	if c == nil {
//		return nil
//	}
//
//	var header http.Header
//	if c.IfMatch != "" {
//		header = make(http.Header)
//		header.Set("if-match", c.IfMatch)
//		return header
//	}
//	return nil
//}
//
//func (c *ReplaceVertexConfig) queryParams() url.Values {
//	if c == nil {
//		return nil
//	}
//
//	var params url.Values
//	if c.WaitForSync != nil {
//		params = make(url.Values)
//		params.Set("waitForSync", strconv.FormatBool(*c.WaitForSync))
//	}
//	return params
//}
//
//func (c *Connection) ReplaceVertex(dbName, graphName, collName, vertexKey string, data interface{}, config *ReplaceVertexConfig) (vertex ReplaceVertexResult, rc int, err error) {
//	path := buildPath(pathConfig{
//		dbName:      dbName,
//		pathFormat:  "/_api/gharial/%s/vertex/%s/%s",
//		pathParams:  []interface{}{graphName, collName, vertexKey},
//		queryParams: config.queryParams(),
//	})
//
//	var body struct {
//		Vertex ReplaceVertexResult `json:"vertex"`
//		Code   int                 `json:"code"`
//	}
//	_, err = c.send(http.MethodPut, path, config.header(), data, &body)
//	if err != nil {
//		return body.Vertex, 0, fmt.Errorf("failed to replace vertex: %v", err)
//	}
//	return body.Vertex, body.Code, nil
//}
//
//type RemoveVertexConfig struct {
//	WaitForSync *bool
//	IfMatch     string
//}
//
//func (c *RemoveVertexConfig) header() http.Header {
//	if c == nil {
//		return nil
//	}
//
//	var header http.Header
//	if c.IfMatch != "" {
//		header = make(http.Header)
//		header.Set("if-match", c.IfMatch)
//		return header
//	}
//	return nil
//}
//
//func (c *RemoveVertexConfig) queryParams() url.Values {
//	if c == nil {
//		return nil
//	}
//
//	var params url.Values
//	if c.WaitForSync != nil {
//		params = make(url.Values)
//		params.Set("waitForSync", strconv.FormatBool(*c.WaitForSync))
//	}
//	return params
//}
//
//func (c *Connection) RemoveVertex(dbName, graphName, collName, vertexKey string, config *RemoveVertexConfig) (removed bool, rc int, err error) {
//	path := buildPath(pathConfig{
//		dbName:      dbName,
//		pathFormat:  "/_api/gharial/%s/vertex/%s/%s",
//		pathParams:  []interface{}{graphName, collName, vertexKey},
//		queryParams: config.queryParams(),
//	})
//
//	var body struct {
//		Removed bool `json:"removed"`
//		Code    int  `json:"code"`
//	}
//	_, err = c.send(http.MethodDelete, path, config.header(), nil, &body)
//	if err != nil {
//		return body.Removed, 0, fmt.Errorf("failed to remove vertex: %v", err)
//	}
//	return body.Removed, body.Code, nil
//}
