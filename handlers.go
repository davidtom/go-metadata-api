package main

import (
	"io/ioutil"
	"net/http"

	"github.com/go-yaml/yaml"
)

type metadata struct {
	Title   string
	Version string
}

// TODO: implement a request logger
func persistMetadata(w http.ResponseWriter, r *http.Request) {
	m := metadata{}

	bs, err := ioutil.ReadAll(r.Body)

	if err != nil {
		logger.Printf("Error reading request body: %v", err)
		http.Error(w, "error reading body", http.StatusBadRequest)

		return
	}

	if err = yaml.Unmarshal(bs, &m); err != nil {
		logger.Printf("error parsing yaml: %v", err)
		http.Error(w, "error parsing yaml", http.StatusBadRequest)

		return
	}

	logger.Printf("successful yaml parse: %+v\n", m)
}

func searchMetadata(w http.ResponseWriter, r *http.Request) {
	logger.Println("success")
}
