package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"google.golang.org/api/iterator"

	"cloud.google.com/go/firestore"
)

func handler(w http.ResponseWriter, r *http.Request) {
	log.Print("Hello world received a request.")
	target := os.Getenv("TARGET")
	if target == "" {
		target = "Bob"
	}
	fmt.Fprintf(w, "Hello %s!\n", target)
}

func connectFirestore() {
	// Sets your Google Cloud Platform project ID.
	projectID := "sd-covid-2"

	// Get a Firestore client.
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	readFromStore(ctx, client)

	// Close client when done.
	defer client.Close()
}

func readFromStore(ctx context.Context, client *firestore.Client) {
	iter := client.Collection("covid_data").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}
		fmt.Printf("Found firestore data at ref: %v", doc.Ref)
	}
}

func main() {
	log.Print("Hello world sample started.")

	http.HandleFunc("/", handler)

	err := os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "./config/sd-covid-2-3c873e023505.json")
	if err != nil {
		log.Println("Error", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	connectFirestore()

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
