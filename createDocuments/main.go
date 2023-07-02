package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
)

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{"http://localhost:8529"},
	})
	if err != nil {
		log.Fatal("Failed to create HTTP connection:", err)
	}

	client, err := driver.NewClient(driver.ClientConfig{
		Connection: conn,
	})
	if err != nil {
		log.Fatal("Failed to create database connection:", err)
	}

	// Create database
	db, err := client.CreateDatabase(context.Background(), "example_users", nil)
	if err != nil {
		log.Fatal("Failed to create database:", err)
	}

	// Create collection
	col, err := db.CreateCollection(context.Background(), "users", nil)
	if err != nil {
		log.Fatal("Failed to create collection:", err)
	}

	// Create documents
	users := []User{
		{Name: "John", Age: 65},
		{Name: "Tina", Age: 25},
		{Name: "George", Age: 31},
	}

	metas, errs, err := col.CreateDocuments(context.Background(), users)
	if err != nil {
		log.Fatal("Failed to create documents:", err)
	} else if err := errs.FirstNonNil(); err != nil {
		log.Fatal("Failed to create documents: first error:", err)
	}

	fmt.Printf("Created documents with keys '%s' in collection '%s' in database '%s'\n", strings.Join(metas.Keys(), ","), col.Name(), db.Name())
}
