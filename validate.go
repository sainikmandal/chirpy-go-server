package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func validateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(errorVal{Error: "Something went wrong"})
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(errorVal{Error: "Chirp is too long"})
		return
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(okVal{Valid: true})
}

type errorVal struct {
	Error string `json:"error"`
}

type okVal struct {
	Valid bool `json:"valid"`
}
