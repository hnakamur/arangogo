package arangogo

import (
	"fmt"
	"net/http"
)

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

type CreateCollectionResult struct {
	ID          string
	Name        string
	WaitForSync bool
	IsVolatile  bool
	IsSystem    bool
	Status      int
	Type        int
}

func (c *Connection) CreateCollection(dbName string, config CreateCollectionConfig) (r CreateCollectionResult, rc int, err error) {
	path := buildPath(pathConfig{
		dbName:     dbName,
		pathFormat: "/_api/collection",
	})

	var body struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		WaitForSync bool   `json:"waitForSync"`
		IsVolatile  bool   `json:"isVolatile"`
		IsSystem    bool   `json:"isSystem"`
		Status      int    `json:"status"`
		Type        int    `json:"type"`
		Code        int    `json:"code"`
	}
	_, err = c.send(http.MethodPost, path, nil, config, &body)
	if err != nil {
		return r, body.Code, fmt.Errorf("failed to create collection: %v", err)
	}
	r = CreateCollectionResult{
		ID:          body.ID,
		Name:        body.Name,
		WaitForSync: body.WaitForSync,
		IsVolatile:  body.IsVolatile,
		IsSystem:    body.IsSystem,
		Status:      body.Status,
		Type:        body.Type,
	}
	return r, body.Code, nil
}

type DropCollectionResult struct {
	ID string
}

func (c *Connection) DropCollection(dbName, collectionName string) (r DropCollectionResult, rc int, err error) {
	path := buildPath(pathConfig{
		dbName:     dbName,
		pathFormat: "/_api/collection/%s",
		pathParams: []interface{}{collectionName},
	})

	var body struct {
		ID   string `json:"id"`
		Code int    `json:"code"`
	}
	_, err = c.send(http.MethodDelete, path, nil, nil, &body)
	if err != nil {
		return r, body.Code, fmt.Errorf("failed to drop collection: %v", err)
	}
	r = DropCollectionResult{
		ID: body.ID,
	}
	return r, body.Code, nil
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

func (c *Connection) TruncateCollection(dbName, collectionName string) (r Collection, rc int, err error) {
	path := buildPath(pathConfig{
		dbName:     dbName,
		pathFormat: "/_api/collection/%s/truncate",
		pathParams: []interface{}{collectionName},
	})

	var body struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		IsSystem bool   `json:"isSystem"`
		Status   int    `json:"status"`
		Type     int    `json:"type"`
		Code     int    `json:"code"`
	}
	_, err = c.send(http.MethodPut, path, nil, nil, &body)
	if err != nil {
		return r, body.Code, fmt.Errorf("failed to truncate collection: %v", err)
	}
	r = Collection{
		ID:       body.ID,
		Name:     body.Name,
		IsSystem: body.IsSystem,
		Status:   body.Status,
		Type:     body.Type,
	}
	return r, body.Code, nil
}
