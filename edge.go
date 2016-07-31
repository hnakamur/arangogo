package arangogo

import (
	"fmt"
	"net/url"
)

type Edge struct {
	ID   string `json:"_id"`
	Key  string `json:"_key"`
	Rev  string `json:"_rev"`
	From string `json:"_from"`
	To   string `json:"_to"`
}

const (
	DirectionIn  = "in"
	DirectionOut = "out"
)

type ListEdgesConfig struct {
	Vertex    string
	Direction string
}

func (c *Connection) ListEdges(dbName, collName string, config *ListEdgesConfig) ([]Edge, error) {
	var body []Edge
	u := dbPrefix(dbName) + "/_api/edges/" + collName
	//var v []string
	//if config != nil {
	//	if config.Vertex != "" {
	//		v = append(v, "vertex="+config.Vertex)
	//	}
	//	if config.Direction != "" {
	//		v = append(v, "direction="+config.Direction)
	//	}
	//}
	//if len(v) > 0 {
	//	u = u + "?" + strings.Join(v, "&")
	//}

	v := url.Values{}
	if config != nil {
		if config.Vertex != "" {
			v.Set("vertex", config.Vertex)
		}
		if config.Direction != "" {
			v.Set("direction", config.Direction)
		}
	}
	if len(v) > 0 {
		u = u + "?" + v.Encode()
	}
	_, err := c.send("GET", u, nil, nil, body)
	if err != nil {
		return nil, fmt.Errorf("failed to list edges: %v", err)
	}
	return body, nil
}
