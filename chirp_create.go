package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/sainikmandal/chirpy-go-server/internal/auth"
	db "github.com/sainikmandal/chirpy-go-server/internal/database"
)

func (cfg *apiConfig) chirpCreateHandler(w http.ResponseWriter, r *http.Request) {
	// Get and validate JWT token
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid authorization header", err)
		return
	}

	// Validate JWT and get user ID
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token", err)
		return
	}

	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	words := strings.Split(params.Body, " ")
	for i, word := range words {
		lowerWord := strings.ToLower(word)
		if _, ok := badWords[lowerWord]; ok {
			words[i] = "****"
		}
	}
	cleanedBody := strings.Join(words, " ")

	chirp, err := cfg.db.CreateChirp(r.Context(), db.CreateChirpParams{
		Body:   cleanedBody,
		UserID: userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, chirp)
}
