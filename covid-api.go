package main

import (
	"fmt"
	"go-covid-api/app"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	absPath, err := filepath.Abs("config/sd-covid-2-3c873e023505.json")
	if err != nil {
		log.Fatalf("Error setting file path: %v", err)
	}

	err = os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", absPath)
	if err != nil {
		log.Fatalf("Error setting env var for firestore credentials: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to http://localhost:%s", port)
	}

	r := app.ConnectFirestore()

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), r))
}
