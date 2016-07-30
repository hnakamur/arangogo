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
			Users: []map[string]interface{}{
				map[string]interface{}{
					"username": "root",
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
	err = c.CreateCollection(dbName, ara.CreateCollectionConfig{Name: collName})
	if err != nil {
		return err
	}

	collections, err := c.ListCollections(dbName)
	if err != nil {
		return err
	}
	for _, c := range collections {
		log.Printf("collection=%v", c)
	}

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

	err = c.TruncateCollection(dbName, collName)
	if err != nil {
		return err
	}

	err = c.DeleteCollection(dbName, collName)
	if err != nil {
		return err
	}

	return nil
}
