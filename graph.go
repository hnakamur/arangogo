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
