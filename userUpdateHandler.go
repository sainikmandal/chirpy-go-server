package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sainikmandal/chirpy-go-server/internal/auth"
	db "github.com/sainikmandal/chirpy-go-server/internal/database"
)

func (cfg *apiConfig) userUpdateHandler(w http.ResponseWriter, r *http.Request) {
	// Get the token from the header
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token", err)
		return
	}

	// Validate the token and get the user ID
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token", err)
		return
	}

	// Parse the request body
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	// Hash the new password
	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
		return
	}

	// Update the user in the database
	user, err := cfg.db.UpdateUser(r.Context(), db.UpdateUserParams{
		ID:             userID,
		Email:          params.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user", err)
		return
	}

	// Return the updated user (without the password)
	respondWithJSON(w, http.StatusOK, struct {
		ID        uuid.UUID `json:"id"`
		Email     string    `json:"email"`
		CreatedAt string    `json:"created_at"`
		UpdatedAt string    `json:"updated_at"`
	}{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	})
}
