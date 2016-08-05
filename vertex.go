package arangogo

import (
	"fmt"
	"net/url"
	"strconv"
)

type CreateVertexConfig struct {
	WaitForSync *bool
}

func (c CreateVertexConfig) urlValues() url.Values {
	if c.WaitForSync != nil {
		return url.Values{
			"waitForSync": []string{
				strconv.FormatBool(*c.WaitForSync),
			},
		}
	}
	return nil
}

func (c *Connection) CreateVertex(dbName, graphName, collName string, data interface{}, config CreateVertexConfig) (DocIDKeyRev, int, error) {
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
	_, err := c.send("POST", path, nil, data, &body)
	if err != nil {
		return body.Vertex, 0, fmt.Errorf("failed to create vertex: %v", err)
	}
	return body.Vertex, body.Code, nil
}
