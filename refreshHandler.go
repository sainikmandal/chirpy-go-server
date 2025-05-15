package main

import (
	"net/http"
	"time"

	"github.com/sainikmandal/chirpy-go-server/internal/auth"
)

func (cfg *apiConfig) refreshHandler(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid authorization header", err)
		return
	}

	userID, err := cfg.db.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid or expired refresh token", err)
		return
	}

	token, err := auth.MakeJWT(userID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, struct {
		Token string `json:"token"`
	}{
		Token: token,
	})
}
