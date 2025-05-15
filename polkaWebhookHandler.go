package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/sainikmandal/chirpy-go-server/internal/auth"
)

type PolkaWebhookRequest struct {
	Event string `json:"event"`
	Data  struct {
		UserID string `json:"user_id"`
	} `json:"data"`
}

func (cfg *apiConfig) polkaWebhookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil || apiKey != cfg.PolkaKey {
		respondWithError(w, http.StatusUnauthorized, "Invalid or missing API key", err)
		return
	}

	var req PolkaWebhookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	if req.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	userID, err := uuid.Parse(req.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	// Check if user exists
	_, err = cfg.db.GetUser(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "User not found", err)
		return
	}

	// Upgrade user to Chirpy Red
	err = cfg.db.UpgradeUserToChirpyRed(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to upgrade user", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
