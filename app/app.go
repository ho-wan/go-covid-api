package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// server - router that wraps around firestore
type server struct {
	router *chi.Mux
	ctx    context.Context
	client *firestore.Client
	data   map[string]interface{}
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

	s.data = s.readFromStore()

	s.routes()

	return r
}

func (s *server) readFromStore() map[string]interface{} {
	client := s.client
	ctx := s.ctx

	// covidData := client.Collection("covid_data")
	allCases := client.Doc("covid_data/all_cases")
	docsnap, err := allCases.Get(ctx)
	if err != nil {
		log.Printf("Failed to get data from firestore: %v", err)
	}

	dataMap := docsnap.Data()
	return dataMap
}

func (s *server) routes() {
	s.router.Get("/", handler)
	s.router.Get("/data", s.handleFetchData())
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "go-covid-api is running!\n")
}

func (s *server) handleFetchData() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO handle data update
		data := s.data

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(data)
	}
}
