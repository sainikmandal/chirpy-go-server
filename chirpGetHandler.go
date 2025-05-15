package main

import (
	"net/http"
	"sort"

	"github.com/google/uuid"
	db "github.com/sainikmandal/chirpy-go-server/internal/database"
)

func (cfg *apiConfig) chirpGetHandler(w http.ResponseWriter, r *http.Request) {
	// Get author_id from query parameters
	authorIDStr := r.URL.Query().Get("author_id")
	var authorID *uuid.UUID
	if authorIDStr != "" {
		id, err := uuid.Parse(authorIDStr)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid author_id", err)
			return
		}
		authorID = &id
	}

	// Get sort parameter from query parameters
	sortOrder := r.URL.Query().Get("sort")
	if sortOrder != "" && sortOrder != "asc" && sortOrder != "desc" {
		respondWithError(w, http.StatusBadRequest, "Invalid sort parameter. Must be 'asc' or 'desc'", nil)
		return
	}
	if sortOrder == "" {
		sortOrder = "asc" // Default to ascending order
	}

	// Get chirps from database
	chirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps", err)
		return
	}

	// Filter by author_id if provided
	if authorID != nil {
		filteredChirps := make([]db.Chirp, 0)
		for _, chirp := range chirps {
			if chirp.UserID == *authorID {
				filteredChirps = append(filteredChirps, chirp)
			}
		}
		chirps = filteredChirps
	}

	// Sort chirps by created_at
	sort.Slice(chirps, func(i, j int) bool {
		if sortOrder == "desc" {
			return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
		}
		return chirps[i].CreatedAt.Before(chirps[j].CreatedAt)
	})

	respondWithJSON(w, http.StatusOK, chirps)
}
