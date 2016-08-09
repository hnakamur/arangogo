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
	ID          string `json:"id"`
	Name        string `json:"name"`
	WaitForSync bool   `json:"waitForSync"`
	IsVolatile  bool   `json:"isVolatile"`
	IsSystem    bool   `json:"isSystem"`
	Status      int    `json:"status"`
	Type        int    `json:"type"`
}

func (c *Connection) CreateCollection(dbName string, config CreateCollectionConfig) (r CreateCollectionResult, rc int, err error) {
	path := buildPath(pathConfig{
		dbName:     dbName,
		pathFormat: "/_api/collection",
	})

	rc, _, err = c.send(http.MethodPost, path, nil, config, &r)
	if err != nil {
		return r, rc, fmt.Errorf("failed to create collection: %v", err)
	}
	return r, rc, nil
}

type DropCollectionResult struct {
	ID string `json:"id"`
}

func (c *Connection) DropCollection(dbName, collectionName string) (r DropCollectionResult, rc int, err error) {
	path := buildPath(pathConfig{
		dbName:     dbName,
		pathFormat: "/_api/collection/%s",
		pathParams: []interface{}{collectionName},
	})

	rc, _, err = c.send(http.MethodDelete, path, nil, nil, &r)
	if err != nil {
		return r, rc, fmt.Errorf("failed to drop collection: %v", err)
	}
	return r, rc, nil
}

func (c *Connection) ListCollections(dbName string) (r []Collection, rc int, err error) {
	path := buildPath(pathConfig{
		dbName:     dbName,
		pathFormat: "/_api/collection",
	})

	var body struct {
		Result []Collection `json:"result"`
	}
	rc, _, err = c.send(http.MethodGet, path, nil, nil, body)
	if err != nil {
		return nil, rc, fmt.Errorf("failed to list collections: %v", err)
	}
	return body.Result, rc, nil
}

func (c *Connection) TruncateCollection(dbName, collectionName string) (r Collection, rc int, err error) {
	path := buildPath(pathConfig{
		dbName:     dbName,
		pathFormat: "/_api/collection/%s/truncate",
		pathParams: []interface{}{collectionName},
	})

	rc, _, err = c.send(http.MethodPut, path, nil, nil, &r)
	if err != nil {
		return r, rc, fmt.Errorf("failed to truncate collection: %v", err)
	}
	return r, rc, nil
}
