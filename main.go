package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
)

var logger = log.New(os.Stdout, "http: ", log.LstdFlags)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}

	r := chi.NewRouter()

	srv := &http.Server{
		Addr:         port,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	r.Route("/v1/metadata", func(r chi.Router) {
		r.Post("/", persistMetadataHandler)
		r.Get("/search", searchMetadataHandler)
	})

	fmt.Println("Server listening on port", port)
	logger.Fatal(srv.ListenAndServe())
}
