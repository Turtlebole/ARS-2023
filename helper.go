package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	ps "github.com/Turtlebole/ARS-2023/poststore"
	"github.com/google/uuid"
)

func decodeBody(r io.Reader) (*ps.Config, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	var cfg ps.Config
	if err := dec.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("error decoding JSON: %s", err)
	}
	return &cfg, nil
}

func decodeGroup(r io.Reader) (*ps.Group, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	var g ps.Group
	if err := dec.Decode(&g); err != nil {
		return nil, fmt.Errorf("error decoding JSON: %s", err)
	}
	return &g, nil
}

func renderJSON(w http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func createId() string {
	return uuid.New().String()
}
