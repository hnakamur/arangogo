package arangogo

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

func (c *Connection) ListGraphs(dbName string, graphsPtr interface{}) (rc int, err error) {
	path := buildPath(pathConfig{
		dbName:     dbName,
		pathFormat: "/_api/gharial",
	})

	var body struct {
		Graphs interface{} `json:"graphs"`
		Code   int         `json:"code"`
	}
	if graphsPtr != nil {
		body.Graphs = graphsPtr
	}
	_, err = c.send(http.MethodGet, path, nil, nil, &body)
	if err != nil {
		return 0, fmt.Errorf("failed to list graphs: %v", err)
	}
	return body.Code, nil
}

type EdgeDefinition struct {
	Collection string   `json:"collection"`
	From       []string `json:"from"`
	To         []string `json:"to"`
}

type CreateGraphConfig struct {
	Name              string           `json:"name"`
	EdgeDefinitions   []EdgeDefinition `json:"edgeDefinitions",omitempty`
	OrphanCollections []string         `json:"orphanCollections",omitempty`
}

type CreateGraphResult struct {
	Name              string           `json:"name"`
	EdgeDefinitions   []EdgeDefinition `json:"edgeDefinitions",omitempty`
	OrphanCollections []string         `json:"orphanCollections",omitempty`
	ID                string           `json:"_id"`
	Rev               string           `json:"_rev"`
}

func (c *Connection) CreateGraph(dbName string, data interface{}) (r CreateGraphResult, rc int, err error) {
	path := buildPath(pathConfig{
		dbName:     dbName,
		pathFormat: "/_api/gharial",
	})

	var body struct {
		Graph CreateGraphResult `json:"graph"`
		Code  int               `json:"code"`
	}
	_, err = c.send(http.MethodPost, path, nil, data, &body)
	if err != nil {
		return body.Graph, body.Code, fmt.Errorf("failed to create graph: %v", err)
	}
	return body.Graph, body.Code, nil
}

type DropGraphConfig struct {
	// NOTE: waitForSync is not documented in query parameters of https://docs.arangodb.com/3.0/HTTP/Gharial/Management.html#drop-a-graph
	WaitForSync     *bool
	DropCollections *bool
	// NOTE: IfMatch is not supported?
}

func (c *DropGraphConfig) queryParams() url.Values {
	if c == nil {
		return nil
	}

	var params url.Values
	if c.WaitForSync != nil || c.DropCollections != nil {
		params = make(url.Values)
	}
	if c.WaitForSync != nil {
		params.Set("waitForSync", strconv.FormatBool(*c.WaitForSync))
	}
	if c.DropCollections != nil {
		params.Set("dropCollections", strconv.FormatBool(*c.DropCollections))
	}
	return params
}

func (c *Connection) DropGraph(dbName, graphName string, config *DropGraphConfig) (removed bool, rc int, err error) {
	path := buildPath(pathConfig{
		dbName:      dbName,
		pathFormat:  "/_api/gharial/%s",
		pathParams:  []interface{}{graphName},
		queryParams: config.queryParams(),
	})

	var body struct {
		Removed bool `json:"removed"`
		Code    int  `json:"code"`
	}
	_, err = c.send(http.MethodDelete, path, nil, nil, &body)
	if err != nil {
		return body.Removed, body.Code, fmt.Errorf("failed to drop edge: %v", err)
	}
	return body.Removed, body.Code, nil
}

func (c *Connection) ListVertexCollections(dbName, graphName string) (collections []string, rc int, err error) {
	path := buildPath(pathConfig{
		dbName:     dbName,
		pathFormat: "/_api/gharial/%s/vertex",
		pathParams: []interface{}{graphName},
	})

	var body struct {
		Collections []string `json:"collections"`
		Code        int      `json:"code"`
	}
	_, err = c.send(http.MethodGet, path, nil, nil, &body)
	if err != nil {
		return nil, body.Code, fmt.Errorf("failed to list vertex collections: %v", err)
	}
	return body.Collections, body.Code, nil
}

type AddVertexCollectionResult struct {
	Name              string           `json:"name"`
	EdgeDefinitions   []EdgeDefinition `json:"edgeDefinitions",omitempty`
	OrphanCollections []string         `json:"orphanCollections",omitempty`
	ID                string           `json:"_id"`
	Rev               string           `json:"_rev"`
}

type AddVertexCollectionConfig struct {
	WaitForSync *bool
}

func (c *AddVertexCollectionConfig) urlValues() url.Values {
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

func (c *Connection) AddVertexCollection(dbName, graphName, collectionName string, config *AddVertexCollectionConfig) (r AddVertexCollectionResult, rc int, err error) {
	path := buildPath(pathConfig{
		dbName:      dbName,
		pathFormat:  "/_api/gharial/%s/vertex",
		pathParams:  []interface{}{graphName},
		queryParams: config.urlValues(),
	})

	payload := struct {
		Collection string `json:"collection"`
	}{
		Collection: collectionName,
	}
	var body struct {
		Graph AddVertexCollectionResult `json:"graph"`
		Code  int                       `json:"code"`
	}
	_, err = c.send(http.MethodPost, path, nil, payload, &body)
	if err != nil {
		return body.Graph, body.Code, fmt.Errorf("failed to add vertex collections: %v", err)
	}
	return body.Graph, body.Code, nil
}

type RemoveVertexCollectionResult struct {
	Name              string           `json:"name"`
	EdgeDefinitions   []EdgeDefinition `json:"edgeDefinitions",omitempty`
	OrphanCollections []string         `json:"orphanCollections",omitempty`
	ID                string           `json:"_id"`
	Rev               string           `json:"_rev"`
}

type RemoveVertexCollectionConfig struct {
	WaitForSync *bool
}

func (c *RemoveVertexCollectionConfig) urlValues() url.Values {
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

func (c *Connection) RemoveVertexCollection(dbName, graphName, collectionName string, config *RemoveVertexCollectionConfig) (r RemoveVertexCollectionResult, rc int, err error) {
	path := buildPath(pathConfig{
		dbName:      dbName,
		pathFormat:  "/_api/gharial/%s/vertex/%s",
		pathParams:  []interface{}{graphName, collectionName},
		queryParams: config.urlValues(),
	})

	var body struct {
		Graph RemoveVertexCollectionResult `json:"graph"`
		Code  int                          `json:"code"`
	}
	_, err = c.send(http.MethodDelete, path, nil, nil, &body)
	if err != nil {
		return body.Graph, body.Code, fmt.Errorf("failed to remove vertex collections: %v", err)
	}
	return body.Graph, body.Code, nil
}

func (c *Connection) ListEdgeDefinitions(dbName, graphName string) (collections []string, rc int, err error) {
	path := buildPath(pathConfig{
		dbName:     dbName,
		pathFormat: "/_api/gharial/%s/edge",
		pathParams: []interface{}{graphName},
	})

	var body struct {
		Collections []string `json:"collections"`
		Code        int      `json:"code"`
	}
	_, err = c.send(http.MethodGet, path, nil, nil, &body)
	if err != nil {
		return nil, body.Code, fmt.Errorf("failed to list edge definitions: %v", err)
	}
	return body.Collections, body.Code, nil
}

type AddEdgeDefinitionResult struct {
	Name              string           `json:"name"`
	EdgeDefinitions   []EdgeDefinition `json:"edgeDefinitions",omitempty`
	OrphanCollections []string         `json:"orphanCollections",omitempty`
	ID                string           `json:"_id"`
	Rev               string           `json:"_rev"`
}

func (c *Connection) AddEdgeDefinition(dbName, graphName string, edgeDefinition interface{}) (r AddEdgeDefinitionResult, rc int, err error) {
	path := buildPath(pathConfig{
		dbName:     dbName,
		pathFormat: "/_api/gharial/%s/edge",
		pathParams: []interface{}{graphName},
	})

	var body struct {
		Graph AddEdgeDefinitionResult `json:"graph"`
		Code  int                     `json:"code"`
	}
	_, err = c.send(http.MethodPost, path, nil, edgeDefinition, &body)
	if err != nil {
		return body.Graph, body.Code, fmt.Errorf("failed to add edge definition: %v", err)
	}
	return body.Graph, body.Code, nil
}

type RemoveEdgeDefinitionResult struct {
	Name              string           `json:"name"`
	EdgeDefinitions   []EdgeDefinition `json:"edgeDefinitions",omitempty`
	OrphanCollections []string         `json:"orphanCollections",omitempty`
	ID                string           `json:"_id"`
	Rev               string           `json:"_rev"`
}

func (c *Connection) RemoveEdgeDefinition(dbName, graphName, definitionName string) (r RemoveEdgeDefinitionResult, rc int, err error) {
	path := buildPath(pathConfig{
		dbName:     dbName,
		pathFormat: "/_api/gharial/%s/edge/%s",
		pathParams: []interface{}{graphName, definitionName},
	})

	var body struct {
		Graph RemoveEdgeDefinitionResult `json:"graph"`
		Code  int                        `json:"code"`
	}
	_, err = c.send(http.MethodDelete, path, nil, nil, &body)
	if err != nil {
		return body.Graph, body.Code, fmt.Errorf("failed to remove edge definition: %v", err)
	}
	return body.Graph, body.Code, nil
}
