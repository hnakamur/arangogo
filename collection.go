package arangogo

import "fmt"

type Collection struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	IsSystem bool   `json:"isSystem"`
	Status   int    `json:"status"`
	Type     int    `json:"type"`
}

func (c *Connection) ListCollections(excludeSystem bool) ([]Collection, error) {
	body := new(struct {
		Result []Collection `json:"result"`
		Error  bool         `json:"error"`
	})
	resp, err := c.send("GET", "/_api/collection", nil, body)
	if err != nil {
		return nil, fmt.Errorf("failed to get collection list: %v", err)
	}
	if body.Error {
		return nil, fmt.Errorf("error in collection list response:%s", string(resp.body))
	}
	var collections []Collection
	if excludeSystem {
		for _, c := range body.Result {
			if !c.IsSystem {
				collections = append(collections, c)
			}
		}
	} else {
		collections = body.Result
	}
	return collections, nil
}
