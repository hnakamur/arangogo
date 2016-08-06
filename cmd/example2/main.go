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

	collName := "startVertices"
	waitForSync := true
	createVertexRes, rc, err := c.CreateVertex(dbName, graphName, collName,
		map[string]string{"name": "Francis"},
		&ara.CreateVertexConfig{WaitForSync: &waitForSync})
	if err != nil {
		return err
	}
	log.Printf("CreateVertex. createVertexRes=%v, rc=%d", createVertexRes, rc)

	modifyVertexRes, rc, err := c.ModifyVertex(dbName, graphName, collName, createVertexRes.Key,
		map[string]interface{}{"age": 26}, &ara.ModifyVertexConfig{WaitForSync: ara.TruePtr()})
	if err != nil {
		return err
	}
	log.Printf("ModifyVertex. modifyVertexRes=%v, rc=%d", modifyVertexRes, rc)

	replaceVertexRes, rc, err := c.ReplaceVertex(dbName, graphName, collName, createVertexRes.Key,
		map[string]interface{}{"name": "Alice Cooper", "age": 26},
		&ara.ReplaceVertexConfig{WaitForSync: ara.TruePtr()})
	if err != nil {
		return err
	}
	log.Printf("ReplaceVertex. replaceVertexRes=%v, rc=%d", replaceVertexRes, rc)

	var vertex interface{}
	rc, err = c.GetVertex(dbName, graphName, collName, createVertexRes.Key, nil, &vertex)
	if err != nil {
		return err
	}
	log.Printf("GetVertex. vertex=%v, rc=%d", vertex, rc)

	createEdgeRes, rc, err := c.CreateEdge(dbName, graphName, collName,
		map[string]interface{}{
			"type":  "friend",
			"_from": "female/alice",
			"_to":   "female/diana",
		}, &ara.CreateEdgeConfig{WaitForSync: ara.TruePtr()})
	if err != nil {
		return err
	}
	log.Printf("CreateEdge. createEdgeRes=%v, rc=%d", createEdgeRes, rc)

	modifyEdgeRes, rc, err := c.ModifyEdge(dbName, graphName, collName, createEdgeRes.Key,
		map[string]interface{}{
			"since": "01.01.2001",
		}, &ara.ModifyEdgeConfig{WaitForSync: ara.TruePtr()})
	if err != nil {
		return err
	}
	log.Printf("ModifyEdge. modifyEdgeRes=%v, rc=%d", modifyEdgeRes, rc)

	replaceEdgeRes, rc, err := c.ReplaceEdge(dbName, graphName, collName, createEdgeRes.Key,
		map[string]interface{}{
			"type":  "divorced",
			"_from": "female/alice",
			"_to":   "male/bob",
		}, &ara.ReplaceEdgeConfig{WaitForSync: ara.TruePtr(), IfMatch: modifyEdgeRes.Rev})
	if err != nil {
		return err
	}
	log.Printf("ReplaceEdge. replaceEdgeRes=%v, rc=%d", replaceEdgeRes, rc)

	var edge interface{}
	rc, err = c.GetEdge(dbName, graphName, collName, createEdgeRes.Key, nil, &edge)
	if err != nil {
		return err
	}
	log.Printf("GetEdge. edge=%v, rc=%d", edge, rc)

	removed, rc, err := c.RemoveVertex(dbName, graphName, collName, createVertexRes.Key,
		&ara.RemoveVertexConfig{WaitForSync: ara.FalsePtr()})
	if err != nil {
		return err
	}
	log.Printf("RemoveVertex. removed=%v, rc=%d", removed, rc)

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
