package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
)

var (
	port string
)

func init() {
	port = os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}
}

func main() {
	r := chi.NewRouter()

	r.Route("/v1/metadata", func(r chi.Router) {
		r.Post("/", persistMetadata)
		r.Get("/search", searchMetadata)
	})

	fmt.Println("Server listening on port", port)
	log.Fatal(http.ListenAndServe(port, r))
}
