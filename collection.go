package arangogo

import "fmt"

type Collection struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	IsSystem bool   `json:"isSystem"`
	Status   int    `json:"status"`
	Type     int    `json:"type"`
}

type CreateCollectionConfig struct {
	JournalSize    int                    `json:"journalSize,omitempty"`
	KeyOptions     map[string]interface{} `json:"keyOptions,omitempty"`
	Name           string                 `json:"name"`
	WaitForSync    *bool                  `json:"waitForSync,omitempty"`
	DoCompact      *bool                  `json:"doCompact,omitempty"`
	IsVolatile     *bool                  `json:"isVolatile,omitempty"`
	ShardKeys      []string               `json:"shardKeys,omitempty"`
	NumberOfShards int                    `json:"numberOfShards,omitempty"`
	IsSystem       *bool                  `json:"isSystem,omitempty"`
	Type           int                    `json:"type,omitempty"`
	IndexBuckets   int                    `json:"indexBuckets,omitempty"`
}

func (c *Connection) CreateCollection(config CreateCollectionConfig) error {
	body := new(struct {
		Error bool `json:"error"`
	})
	resp, err := c.send("POST", "/_api/collection", config, body)
	if err != nil {
		return fmt.Errorf("failed to create collection: %v", err)
	}
	if body.Error {
		return fmt.Errorf("error in create collection response:%s", string(resp.body))
	}
	return nil
}

func (c *Connection) DeleteCollection(name string) error {
	body := new(struct {
		Error bool `json:"error"`
	})
	resp, err := c.send("DELETE", "/_api/collection/"+name, nil, body)
	if err != nil {
		return fmt.Errorf("failed to delete collection: %v", err)
	}
	if body.Error {
		return fmt.Errorf("error in delete collection response:%s", string(resp.body))
	}
	return nil
}

func (c *Connection) ListCollections() ([]Collection, error) {
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
	return body.Result, nil
}

func (c *Connection) TruncateCollection(name string) error {
	body := new(struct {
		Error bool `json:"error"`
	})
	resp, err := c.send("PUT", "/_api/collection/"+name+"/truncate", nil, body)
	if err != nil {
		return fmt.Errorf("failed to truncate collection: %v", err)
	}
	if body.Error {
		return fmt.Errorf("error in truncate collection response:%s", string(resp.body))
	}
	return nil
}
