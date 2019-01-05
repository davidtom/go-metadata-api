package main

import (
	"errors"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/go-validator/validator"
	"github.com/go-yaml/yaml"
)

type metadata struct {
	Title       string `yaml:"title" validate:"nonzero"`
	Version     string `yaml:"version" validate:"nonzero"`
	Maintainers []struct {
		Name  string `yaml:"name" validate:"nonzero"`
		Email string `yaml:"email" validate:"nonzero"`
	} `yaml:"maintainers"`
	Company     string `yaml:"company" validate:"nonzero"`
	Website     string `yaml:"website" validate:"nonzero"`
	Source      string `yaml:"source" validate:"nonzero"`
	License     string `yaml:"license" validate:"nonzero"`
	Description string `yaml:"description"`
}

var emailRegexp = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// TODO: implement a request logger

// Handlers
func persistMetadataHandler(w http.ResponseWriter, r *http.Request) {
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

	if err = validateMetadata(&m); err != nil {
		logger.Printf("invalid yaml: %v", err)
		rejectionReason := "invalid yaml " + err.Error()
		http.Error(w, rejectionReason, http.StatusBadRequest)

		return
	}

	storage.set(m.Title, &m)
	logger.Printf("successful yaml upload: %+v\n", m)

	w.WriteHeader(http.StatusNoContent)
}

func searchMetadataHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	logger.Printf("query: %+v", q)

	results := storage.get()

	w.Header().Set("Content-Type", "application/x-yaml")
	e := yaml.NewEncoder(w)
	defer e.Close()

	for _, yaml := range results {
		e.Encode(yaml)
	}
}

// Helper Methods
func validateMetadata(m *metadata) error {
	if err := validator.Validate(m); err != nil {
		return err
	}

	for _, maintainer := range m.Maintainers {
		// validate email syntax only in the interest of speed
		if validEmail := emailRegexp.MatchString(maintainer.Email); validEmail != true {
			return errors.New("invalid email")
		}
	}

	return nil
}
