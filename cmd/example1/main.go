package main

import (
	"log"

	ara "github.com/hnakamur/arangogo"
)

func main() {
	err := run("root", "root")
	if err != nil {
		panic(err)
	}
}

func contains(array []string, elem string) bool {
	for _, item := range array {
		if item == elem {
			return true
		}
	}
	return false
}

func run(username, password string) (err error) {
	c, err := ara.NewConnection(&ara.Config{Username: username, Password: password})
	if err != nil {
		return err
	}

	databases, err := c.ListDatabases()
	if err != nil {
		return err
	}
	log.Printf("databases=%v", databases)

	dbName := "foo"
	if !contains(databases, dbName) {
		err = c.CreateDatabase(ara.CreateDatabaseConfig{
			Name: dbName,
			Users: []ara.CreateDatabaseConfigUser{
				{
					Username: "root",
				},
			},
		})
		if err != nil {
			return err
		}
	}
	defer func() {
		err = c.DropDatabase(dbName)
		if err != nil {
			return
		}

		var userDatabases []string
		userDatabases, err = c.ListUserDatabases()
		if err != nil {
			return
		}
		log.Printf("userDatabases=%v", userDatabases)
	}()

	collName := "mycollection"
	createCollectionRes, rc, err := c.CreateCollection(dbName, ara.CreateCollectionConfig{Name: collName})
	if err != nil {
		return err
	}
	log.Printf("CreateCollection. res=%v, rc=%d", createCollectionRes, rc)

	collections, rc, err := c.ListCollections(dbName)
	if err != nil {
		return err
	}
	log.Printf("ListCollections. res=%v, rc=%d", collections, rc)

	data := map[string]interface{}{
		"name": "Alice",
	}
	var docBody struct {
		ID   string `json:"_id"`
		Key  string `json:"_key"`
		Rev  string `json:"_rev"`
		Name string `json:"name"`
	}
	doc, rc, err := c.CreateDocument(dbName, collName, data, &ara.CreateDocumentConfig{ReturnNew: ara.TruePtr()}, &docBody)
	if err != nil {
		return err
	}
	log.Printf("CreateDocument. doc=%v, rc=%v, docBody=%v", doc, rc, docBody)

	data2 := []map[string]interface{}{
		{"name": "Alice"},
		{"name": "Bob"},
		{"name": "Charlie"},
	}
	//data2 := `{1:"Foo"},{2:"Bad"}`
	//data2 := `[{"name":"Foo"},{"name":"Bar"}]`
	docs, rc, err := c.CreateDocuments(dbName, collName, data2, nil)
	if err != nil {
		log.Printf("err=%v", err)
		return err
	}
	log.Printf("CreateDocuments. docs=%v, rc=%d", docs, rc)

	var docBody2 struct {
		ID   string `json:"_id"`
		Key  string `json:"_key"`
		Rev  string `json:"_rev"`
		Name string `json:"name"`
	}
	doc, rc, err = c.RemoveDocument(dbName, collName, doc.Key, &ara.RemoveDocumentConfig{IfMatch: doc.Rev, ReturnOld: ara.TruePtr()}, &docBody2)
	if err != nil {
		log.Printf("err=%v", err)
		return err
	}
	log.Printf("RemoveDocument. doc=%v, rc=%v, docBody=%v", doc, rc, docBody2)

	edgeCollName := "myedges"
	createCollectionRes, rc, err = c.CreateCollection(dbName, ara.CreateCollectionConfig{Name: edgeCollName})
	if err != nil {
		return err
	}
	log.Printf("CreateCollection. res=%v, rc=%d", createCollectionRes, rc)
	edges := []map[string]interface{}{
		{
			"label": docs[0].Key + "->" + docs[1].Key,
			"_from": docs[0].ID,
			"_to":   docs[1].ID,
		},
		{
			"label": docs[0].Key + "->" + docs[2].Key,
			"_from": docs[0].ID,
			"_to":   docs[2].ID,
		},
	}
	waitForSync := true
	edgeDocs, rc, err := c.CreateDocuments(dbName, edgeCollName, edges, &ara.CreateDocumentsConfig{
		WaitForSync: &waitForSync,
	})
	if err != nil {
		log.Printf("err=%v", err)
		return err
	}
	log.Printf("created edge documents=%v, rc=%d", edgeDocs, rc)

	collection, rc, err := c.TruncateCollection(dbName, collName)
	if err != nil {
		return err
	}
	log.Printf("TruncateCollection. res=%v, rc=%d", collection, rc)

	dropCollectionRes, rc, err := c.DropCollection(dbName, collName)
	if err != nil {
		return err
	}
	log.Printf("DropCollection. res=%v, rc=%d", dropCollectionRes, rc)

	return nil
}
