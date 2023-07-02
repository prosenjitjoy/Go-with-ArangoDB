package main

import (
	"context"
	"fmt"
	"log"

	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
)

type Book struct {
	Title   string
	NoPages int
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

	// Create "example_books" database
	db, err := client.CreateDatabase(context.Background(), "example_books", nil)
	if err != nil {
		log.Fatal("Failed to create database:", err)
	}

	// Create "books" collection
	col, err := db.CreateCollection(context.Background(), "books", nil)
	if err != nil {
		log.Fatal("Failed to create collection:", err)
	}

	// Create document
	book := Book{
		Title:   "Arangodb Cookbook",
		NoPages: 257,
	}
	meta, err := col.CreateDocument(context.Background(), book)
	if err != nil {
		log.Fatal("Failed to create document:", err)
	}
	fmt.Printf("Created document in collection '%s' in database '%s'\n", col.Name(), db.Name())

	// Read the document back
	var result Book
	_, err = col.ReadDocument(context.Background(), meta.Key, &result)
	if err != nil {
		log.Fatal("Failed to read document:", err)
	}
	fmt.Printf("Read book '%+v'\n", result)
}
