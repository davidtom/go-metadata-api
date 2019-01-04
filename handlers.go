package main

import (
	"io/ioutil"
	"log"
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
	defer r.Body.Close()

	if err != nil {
		msg := "error parsing request body"

		log.Println(msg, err)
		http.Error(w, msg, 500)

		return
	}

	if err = yaml.Unmarshal(bs, &m); err != nil {
		msg := "error parsing yaml"

		log.Println(msg, err)
		http.Error(w, msg, 500)

		return
	}

	log.Printf("successful yaml parse: %+v\n", m)
}

func searchMetadata(w http.ResponseWriter, r *http.Request) {
	log.Println("success")
}
