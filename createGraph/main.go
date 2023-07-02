package main

import (
	"context"
	"log"

	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
)

type MyObject struct {
	Name string `json:"_key"`
	Age  int    `json:"age"`
}

type MyEdgeObject struct {
	From string `json:"_from"`
	To   string `json:"_to"`
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
	})
	if err != nil {
		log.Fatal("Failed to create database connection:", err)
	}

	// Create database
	db, err := client.CreateDatabase(ctx, "my_graph_db", nil)
	if err != nil {
		log.Fatal("Failed to create database:", err)
	}

	// define the edgeCollection to store the edges
	var edgeDefinition driver.EdgeDefinition
	edgeDefinition.Collection = "myEdgeCollection"

	// define a set of collections where an edge is going out
	edgeDefinition.From = []string{"myCollection1", "myCollection2"}

	// repeat this for the collections where an edge is going into
	edgeDefinition.To = []string{"myCollection1", "myCollection3"}

	// A graph can contain additional vertex collections, defined in the set of orphan collections
	var options driver.CreateGraphOptions
	options.OrphanVertexCollections = []string{"myCollection4", "myCollection5"}
	options.EdgeDefinitions = []driver.EdgeDefinition{edgeDefinition}

	// now it's possible to create a graph
	graph, err := db.CreateGraphV2(ctx, "myGraph", &options)
	if err != nil {
		log.Fatal("Failed to create graph:", err)
	}

	// add vertex
	vertexCollection1, err := graph.VertexCollection(ctx, "myCollection1")
	if err != nil {
		log.Fatal("Failed to get vertex collection:", err)
	}

	myObjects := []MyObject{
		{Name: "Homer", Age: 38},
		{Name: "Marge", Age: 36},
	}

	_, _, err = vertexCollection1.CreateDocuments(ctx, myObjects)
	if err != nil {
		log.Fatal("Failed to create vertex documents:", err)
	}

	// add edge
	edgeCollection, _, err := graph.EdgeCollection(ctx, "myEdgeCollection")
	if err != nil {
		log.Fatal("Failed to select edge collection:", err)
	}

	edge := MyEdgeObject{
		From: "myCollection1/Homer",
		To:   "myCollection1/Marge",
	}
	_, err = edgeCollection.CreateDocument(ctx, edge)
	if err != nil {
		log.Fatal("Failed to create edge document:", err)
	}

	// delete graph
	graph.Remove(ctx)
}
