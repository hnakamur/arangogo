package main

import (
	"log"

	"github.com/hnakamur/arangogo"
)

func main() {
	err := run("root", "root")
	if err != nil {
		panic(err)
	}
}

func run(username, password string) error {
	c, err := arangogo.NewConnection(&arangogo.Config{Username: username, Password: password})
	if err != nil {
		return err
	}

	name := "foo"
	err = c.CreateDatabase(name, []interface{}{
		map[string]interface{}{
			"username": "root",
		},
	})
	if err != nil {
		return err
	}

	databases, err := c.ListDatabases()
	if err != nil {
		return err
	}
	log.Printf("databases=%v", databases)

	err = c.DropDatabase(name)
	if err != nil {
		return err
	}

	userDatabases, err := c.ListUserDatabases()
	if err != nil {
		return err
	}
	log.Printf("userDatabases=%v", userDatabases)

	collections, err := c.ListCollections(false)
	if err != nil {
		return err
	}
	for _, c := range collections {
		log.Printf("collection=%v", c)
	}

	return nil
}
