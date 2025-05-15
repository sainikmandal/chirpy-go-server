package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/sainikmandal/chirpy-go-server/internal/auth"
)

func (cfg *apiConfig) chirpDeleteHandler(w http.ResponseWriter, r *http.Request) {
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

	// Get the chirp ID from the URL
	chirpIDStr := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	// Get the chirp to check ownership
	chirp, err := cfg.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found", err)
		return
	}

	// Check if the user is the author of the chirp
	if chirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "Not authorized to delete this chirp", nil)
		return
	}

	// Delete the chirp
	err = cfg.db.DeleteChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete chirp", err)
		return
	}

	// Return 204 No Content
	w.WriteHeader(http.StatusNoContent)
}
