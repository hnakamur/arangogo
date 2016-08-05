package arangogo

import "fmt"

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

func (c *Connection) CreateGraph(dbName string, config CreateGraphConfig) (interface{}, error) {
	var body interface{}
	u := dbPrefix(dbName) + "/_api/gharial"
	_, err := c.send("POST", u, nil, config, &body)
	if err != nil {
		return nil, fmt.Errorf("failed to create graph: %v", err)
	}
	return body, nil
}

func (c *Connection) ListGraphs(dbName string) ([]interface{}, error) {
	var body struct {
		Graphs []interface{} `json:"graphs"`
	}
	u := dbPrefix(dbName) + "/_api/gharial"
	_, err := c.send("GET", u, nil, nil, &body)
	if err != nil {
		return nil, fmt.Errorf("failed to create graph: %v", err)
	}
	return body.Graphs, nil
}

type DropGraphConfig struct {
	Name            string
	DropCollections bool
}

func (c *Connection) DropGraph(dbName string, config DropGraphConfig) error {
	u := dbPrefix(dbName) + "/_api/gharial/" + config.Name
	if config.DropCollections {
		u += "?dropCollections=true"
	}
	_, err := c.send("DELETE", u, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to delete  graph: %v", err)
	}
	return nil
}

func (c *Connection) ListVertexCollections(dbName, graphName string) ([]string, error) {
	var body struct {
		Collections []string `json:"collections"`
	}
	u := dbPrefix(dbName) + "/_api/gharial/" + graphName + "/vertex"
	_, err := c.send("GET", u, nil, nil, &body)
	if err != nil {
		return nil, fmt.Errorf("failed to list vertex collections: %v", err)
	}
	return body.Collections, nil
}

type Graph struct {
	Name              string           `json:"name"`
	EdgeDefinitions   []EdgeDefinition `json:"edgeDefinitions",omitempty`
	OrphanCollections []string         `json:"orphanCollections",omitempty`
	ID                string           `json:"_id"`
	Rev               string           `json:"_rev"`
}

func (c *Connection) AddVertexCollection(dbName, graphName, collectionName string) (Graph, error) {
	payload := struct {
		Collection string `json:"collection"`
	}{
		Collection: collectionName,
	}
	var body struct {
		Graph Graph `json:"graph"`
	}
	u := dbPrefix(dbName) + "/_api/gharial/" + graphName + "/vertex"
	_, err := c.send("POST", u, nil, payload, &body)
	if err != nil {
		return body.Graph, fmt.Errorf("failed to add vertex collections: %v", err)
	}
	return body.Graph, nil
}

func (c *Connection) RemoveVertexCollection(dbName, graphName, collectionName string) (Graph, error) {
	var body struct {
		Graph Graph `json:"graph"`
	}
	u := dbPrefix(dbName) + "/_api/gharial/" + graphName + "/vertex/" + collectionName
	_, err := c.send("DELETE", u, nil, nil, &body)
	if err != nil {
		return body.Graph, fmt.Errorf("failed to remove vertex collections: %v", err)
	}
	return body.Graph, nil
}

func (c *Connection) ListEdgeDefinitions(dbName, graphName string) ([]string, error) {
	var body struct {
		Collections []string `json:"collections"`
	}
	u := dbPrefix(dbName) + "/_api/gharial/" + graphName + "/edge"
	_, err := c.send("GET", u, nil, nil, &body)
	if err != nil {
		return nil, fmt.Errorf("failed to list edge definitions: %v", err)
	}
	return body.Collections, nil
}

func (c *Connection) AddEdgeDefinition(dbName, graphName string, edgeDefinition EdgeDefinition) (Graph, error) {
	var body struct {
		Graph Graph `json:"graph"`
	}
	u := dbPrefix(dbName) + "/_api/gharial/" + graphName + "/edge"
	_, err := c.send("POST", u, nil, edgeDefinition, &body)
	if err != nil {
		return body.Graph, fmt.Errorf("failed to add edge definition: %v", err)
	}
	return body.Graph, nil
}

func (c *Connection) RemoveEdgeDefinition(dbName, graphName, definitionName string) (Graph, error) {
	var body struct {
		Graph Graph `json:"graph"`
	}
	u := dbPrefix(dbName) + "/_api/gharial/" + graphName + "/edge/" + definitionName
	_, err := c.send("DELETE", u, nil, nil, &body)
	if err != nil {
		return body.Graph, fmt.Errorf("failed to remove edge definition: %v", err)
	}
	return body.Graph, nil
}
