package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ho-wan/go-covid-api/app"

	"google.golang.org/api/iterator"

	"cloud.google.com/go/firestore"
)

func main() {
	// http.HandleFunc("/", handler)

	// err := os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "./config/sd-covid-2-3c873e023505.json")
	// if err != nil {
	// 	log.Println("Error", err)
	// }

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to http://localhost:%s", port)
	}

	srv := app.connectFirestore()
	srv.routes()

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
