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
	d, err := arangogo.NewDatabase(&arangogo.Config{Username: username, Password: password})
	if err != nil {
		return err
	}

	name := "foo"
	err = d.CreateDatabase(name, []interface{}{
		map[string]interface{}{
			"username": "root",
		},
	})
	if err != nil {
		return err
	}

	databases, err := d.ListDatabases()
	if err != nil {
		return err
	}
	log.Printf("databases=%v", databases)

	err = d.DropDatabase(name)
	if err != nil {
		return err
	}

	userDatabases, err := d.ListUserDatabases()
	if err != nil {
		return err
	}
	log.Printf("userDatabases=%v", userDatabases)
	return nil
}
