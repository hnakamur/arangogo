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

	graphs, err := c.ListGraphs(dbName)
	if err != nil {
		return err
	}
	for i, graph := range graphs {
		log.Printf("i=%d, graph=%v", i, graph)
	}

	collections, err := c.ListVertexCollections(dbName, graphName)
	if err != nil {
		return err
	}
	for i, collection := range collections {
		log.Printf("ListVertexCollections. i=%d, collection=%s", i, collection)
	}

	graph2, err := c.AddVertexCollection(dbName, graphName, "otherVertices")
	if err != nil {
		return err
	}
	log.Printf("AddVertexCollection. graph=%v", graph2)

	graph2, err = c.AddEdgeDefinition(dbName, graphName, ara.EdgeDefinition{
		Collection: "works_in",
		From: []string{
			"female",
			"male",
		},
		To: []string{
			"city",
		},
	})
	if err != nil {
		return err
	}
	log.Printf("AddEdgeDefinition. graph=%v", graph2)

	collections, err = c.ListEdgeDefinitions(dbName, graphName)
	if err != nil {
		return err
	}
	for i, collection := range collections {
		log.Printf("ListEdgeDefinitions. i=%d, collection=%s", i, collection)
	}

	graph2, err = c.RemoveEdgeDefinition(dbName, graphName, "works_in")
	if err != nil {
		return err
	}
	log.Printf("RemoveEdgeDefinition. graph=%v", graph2)

	graph2, err = c.RemoveVertexCollection(dbName, graphName, "otherVertices")
	if err != nil {
		return err
	}
	log.Printf("RemoveVertexCollection. graph=%v", graph2)

	err = c.DropGraph(dbName, ara.DropGraphConfig{Name: graphName})
	if err != nil {
		return err
	}

	graphs, err = c.ListGraphs(dbName)
	if err != nil {
		return err
	}
	log.Printf("len(graphs)=%d", len(graphs))

	return nil
}
