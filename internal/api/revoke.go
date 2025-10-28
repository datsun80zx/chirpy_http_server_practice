package api

import (
	"net/http"

	"github.com/datsun80zx/chirpy_http_server_practice.git/internal/auth"
)

func (h *Handler) Revoke(w http.ResponseWriter, r *http.Request) {
	// Get refresh token from Authorization header
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "missing or invalid token", err)
		return
	}

	// Revoke the refresh token
	err = h.Config.Database.RevokeRefreshToken(r.Context(), refreshToken)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "problem revoking token", err)
		return
	}

	// Respond with 204 No Content
	w.WriteHeader(http.StatusNoContent)
}
