package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"google.golang.org/api/iterator"
)

// server - router that wraps around firestore
type server struct {
	router *chi.Mux
	ctx    context.Context
	client *firestore.Client
}

// ConnectFirestore - connects to firestore and returns Handler
func ConnectFirestore() http.Handler {
	// Sets your Google Cloud Platform project ID.
	projectID := "sd-covid-2"

	err := os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "./config/sd-covid-2-3c873e023505.json")
	if err != nil {
		log.Println("Error", err)
	}

	// Get a Firestore client.
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	r := chi.NewRouter()
	r.Use(
		middleware.Logger,
	)

	srv := &server{
		ctx:    ctx,
		client: client,
		router: r,
	}

	srv.routes()

	readFromStore(ctx, client)

	// Close client when done.
	defer client.Close()

	return r
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
		fmt.Printf("Found firestore data at ref: %v\n", doc.Ref)
	}
}

func (s *server) routes() {
	s.router.Get("/", handler)
	s.router.Get("/about", s.handleAbout())
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Print("Hello world received a request.")
	target := os.Getenv("TARGET")
	if target == "" {
		target = "Bob"
	}
	fmt.Fprintf(w, "Hello %s!\n", target)
}

func (s *server) handleAbout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "this is handleAbout")
		// use thing
	}
}
