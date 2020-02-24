package main

import (
	"fmt"
	"go-covid-api/app"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to http://localhost:%s", port)
	}

	r := app.ConnectFirestore()

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), r))
}
