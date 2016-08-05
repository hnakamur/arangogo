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

func (c *Connection) CreateCollection(dbName string, config CreateCollectionConfig) error {
	_, err := c.send("POST", dbPrefix(dbName)+"/_api/collection", nil, config, nil)
	if err != nil {
		return fmt.Errorf("failed to create collection: %v", err)
	}
	return nil
}

func (c *Connection) DeleteCollection(dbName string, name string) error {
	_, err := c.send("DELETE", dbPrefix(dbName)+"/_api/collection/"+name, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to delete collection: %v", err)
	}
	return nil
}

func (c *Connection) ListCollections(dbName string) ([]Collection, error) {
	body := new(struct {
		Result []Collection `json:"result"`
	})
	_, err := c.send("GET", dbPrefix(dbName)+"/_api/collection", nil, nil, body)
	if err != nil {
		return nil, fmt.Errorf("failed to get collection list: %v", err)
	}
	return body.Result, nil
}

func (c *Connection) TruncateCollection(dbName string, name string) error {
	_, err := c.send("PUT", dbPrefix(dbName)+"/_api/collection/"+name+"/truncate", nil, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to truncate collection: %v", err)
	}
	return nil
}
