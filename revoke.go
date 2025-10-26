package main

import (
	"net/http"

	"github.com/datsun80zx/chirpy_http_server_practice.git/internal/auth"
)

func (cfg *apiConfig) Revoke(w http.ResponseWriter, r *http.Request) {
	// Get refresh token from Authorization header
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "missing or invalid token", err)
		return
	}

	// Revoke the refresh token
	err = cfg.database.RevokeRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "problem revoking token", err)
		return
	}

	// Respond with 204 No Content
	w.WriteHeader(http.StatusNoContent)
}
