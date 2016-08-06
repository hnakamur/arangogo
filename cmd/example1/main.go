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
	doc, err := c.CreateDocument(dbName, collName, data, nil)
	if err != nil {
		return err
	}
	log.Printf("created document=%v", *doc)

	data2 := []map[string]interface{}{
		{"name": "Alice"},
		{"name": "Bob"},
		{"name": "Charlie"},
	}
	//data2 := `{1:"Foo"},{2:"Bad"}`
	//data2 := `[{"name":"Foo"},{"name":"Bar"}]`
	docs, err := c.CreateDocuments(dbName, collName, data2, nil)
	if err != nil {
		log.Printf("err=%v", err)
		return err
	}
	log.Printf("created documents=%v", docs)

	err = c.DeleteDocument(dbName, collName, doc.Key, &ara.DeleteDocumentConfig{IfMatch: doc.Rev})
	if err != nil {
		log.Printf("err=%v", err)
		return err
	}

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
	edgeDocs, err := c.CreateDocuments(dbName, edgeCollName, edges, &ara.CreateDocumentConfig{
		WaitForSync: &waitForSync,
	})
	if err != nil {
		log.Printf("err=%v", err)
		return err
	}
	log.Printf("created edge documents=%v", edgeDocs)

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
