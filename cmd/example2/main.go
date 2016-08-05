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

	dbName := "mygraphdb1"
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

	graphName := "myGraph"
	graph, err := c.CreateGraph(dbName, ara.CreateGraphConfig{
		Name: graphName,
		EdgeDefinitions: []ara.EdgeDefinition{
			{
				Collection: "edges",
				From: []string{
					"startVertices",
				},
				To: []string{
					"endVertices",
				},
			},
		},
	})
	if err != nil {
		return err
	}
	log.Printf("graph=%v", graph)

	return nil
}
