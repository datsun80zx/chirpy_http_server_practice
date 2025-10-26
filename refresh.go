package main

import (
	"net/http"
	"time"

	"github.com/datsun80zx/chirpy_http_server_practice.git/internal/auth"
)

func (cfg *apiConfig) Refresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	// Get refresh token from Authorization header
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "missing or invalid token", err)
		return
	}

	// Get user from refresh token (this also validates the token)
	user, err := cfg.database.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid or expired refresh token", err)
		return
	}

	// Create new access token that expires in 1 hour
	accessToken, err := auth.MakeJWT(user.ID, cfg.tokenString, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "problem creating token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		Token: accessToken,
	})
}
