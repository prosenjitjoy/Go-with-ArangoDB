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
	ctx := context.Background()

	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{"http://localhost:8529"},
	})
	if err != nil {
		log.Fatal("Failed to create HTTP connection:", err)
	}

	client, err := driver.NewClient(driver.ClientConfig{
		Connection: conn,
		// Authentication: driver.BasicAuthentication("root", "password"),
	})
	if err != nil {
		log.Fatal("Failed to create database connection:", err)
	}

	var db driver.Database

	// check if database already exists
	db_exists, err := client.DatabaseExists(ctx, "example")
	if err != nil {
		log.Fatal("Failed to check if database exists:", err)
	}

	if db_exists {
		fmt.Println("That db exists already")
		db, err = client.Database(ctx, "example")
		if err != nil {
			log.Fatal("Failed to create database:", err)
		}
	} else {
		db, err = client.CreateDatabase(ctx, "example", nil)
		if err != nil {
			log.Fatal("Failed to create database:", err)
		}
	}

	// check if collection already exists
	col_exists, err := db.CollectionExists(ctx, "users")
	if err != nil {
		log.Fatal("Failed to check if database exists:", err)
	}

	if col_exists {
		fmt.Println("That collection exists already")
		PrintCollection(ctx, db, "users")
	} else {
		col, err := db.CreateCollection(ctx, "users", nil)
		if err != nil {
			log.Fatal("Failed to create collection:", err)
		}

		// Create documents
		users := []User{
			{Name: "John", Age: 65},
			{Name: "Tina", Age: 25},
			{Name: "George", Age: 31},
		}
		metas, errs, err := col.CreateDocuments(ctx, users)
		if err != nil {
			log.Fatal("Failed to create document:", err)
		} else if err := errs.FirstNonNil(); err != nil {
			log.Fatal("Failed to create documents: first error:", err)
		}

		fmt.Printf("Create documents with keys '%s' in collection '%s' in database '%s'\n", strings.Join(metas.Keys(), ","), col.Name(), db.Name())
	}
}

func PrintCollection(ctx context.Context, db driver.Database, name string) {
	queryString := fmt.Sprintf("FOR doc IN %s LIMIT 10 RETURN doc", name)
	cursor, err := db.Query(ctx, queryString, nil)
	if err != nil {
		log.Fatal("Failed to execute query:", err)
	}
	defer cursor.Close()

	for {
		var doc User
		meta, err := cursor.ReadDocument(ctx, &doc)

		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			log.Fatal("Failed to read document:", err)
		}

		fmt.Println("Dot doc", meta, doc)
	}
}

// Query documents, fetching the total count
func demoCount(db driver.Database) {
	ctx := driver.WithQueryCount(context.Background())
	query := "FOR d IN myCollection RETURN d"
	cursor, err := db.Query(ctx, query, nil)
	if err != nil {
		log.Fatal("Failed to execute query:", err)
	}
	defer cursor.Close()
	fmt.Printf("Qeury yields %d documents\n", cursor.Count())
}

// Query documents, with bind variables
func demoBind(db driver.Database) {
	ctx := driver.WithQueryCount(context.Background())
	query := "FOR d IN myCollection FILTER d.name == @myVar RETURN d"
	bindVars := map[string]interface{}{
		"myVar": "Some name",
	}
	cursor, err := db.Query(ctx, query, bindVars)
	if err != nil {
		log.Fatal("Failed to execute query:", err)
	}
	defer cursor.Close()
	fmt.Printf("Query yields %d documents\n", cursor.Count())
}
