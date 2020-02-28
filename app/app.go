package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// server - router that wraps around firestore
type server struct {
	router *chi.Mux
	data   map[string]interface{}
}

// ConnectFirestore - connects to firestore and returns Handler
func ConnectFirestore() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	s := &server{
		router: r,
		data:   readDataFromStore(),
	}
	s.routes()

	return r
}

func readDataFromStore() map[string]interface{} {
	// TODO move project-id to dotenv
	projectID := "sd-covid-2"

	// Get a Firestore client.
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	// Close client when done (need a new client to read again)
	defer client.Close()

	allCases := client.Doc("covid_data/all_cases")
	docsnap, err := allCases.Get(ctx)
	if err != nil {
		log.Printf("Failed to get data from firestore: %v", err)
	}

	fmt.Printf("Reading firestore data: %v\n", docsnap)
	return docsnap.Data()
}

func (s *server) routes() {
	s.router.Get("/", handler)
	s.router.Get("/data", s.handleFetchData())
	s.router.Get("/update", s.updateDataHandler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "go-covid-api is running!\n")
}

// updateDataHandler - updates data by fetching from firestore
func (s *server) updateDataHandler(w http.ResponseWriter, r *http.Request) {
	s.data = readDataFromStore()
	fmt.Fprintf(w, "data updated!\n")
}

// handleFetchData - returns data in state as json
func (s *server) handleFetchData() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(s.data)
	}
}
