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
	// TODO move project-id and path to json to dotenv
	projectID := "sd-covid-2"
	err := os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "./config/sd-covid-2-3c873e023505.json")
	if err != nil {
		log.Fatalf("Error setting env var for firestore credentials: %v", err)
	}

	// Get a Firestore client.
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	// Close client when done.
	defer client.Close()

	r := chi.NewRouter()
	r.Use(
		middleware.Logger,
	)

	s := &server{
		ctx:    ctx,
		client: client,
		router: r,
	}

	s.routes()

	return r
}

func (s *server) readFromStore() map[string]interface{} {
	client := s.client
	ctx := s.ctx

	iter := client.Collection("covid_data").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate over data: %v", err)
		}
		fmt.Printf("Reading firestore data: %v\n", doc)
		return doc.Data()
	}
	return nil
}

func (s *server) routes() {
	s.router.Get("/", handler)
	s.router.Get("/data", s.handleData())
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "go-covid-api is running!\n")
}

func (s *server) handleData() http.HandlerFunc {
	data := s.readFromStore()

	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "data: %v\n", data)
	}
}
