package main

import (
	"encoding/json"
	"net/http"

	"github.com/sainikmandal/chirpy-go-server/internal/auth"
	db "github.com/sainikmandal/chirpy-go-server/internal/database"
)

func (cfg *apiConfig) userCreateHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	// Hash the password
	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
		return
	}

	// Create the user in the database
	user, err := cfg.db.CreateUser(r.Context(), db.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}

	// Return the created user using UserResponse type
	respondWithJSON(w, http.StatusCreated, UserResponse{
		ID:          user.ID.String(),
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	})
}
