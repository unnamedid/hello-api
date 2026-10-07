package main

import (
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	environment := os.Getenv("ENVIRONMENT")
	if environment == "" {
		environment = "development"
	}
	srv := &http.Server{
		Addr:              ":8080",
		Handler:           newHandler(environment),
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Printf("hello-api listening on %s (environment=%s)", srv.Addr, environment)
	log.Fatal(srv.ListenAndServe())
}
