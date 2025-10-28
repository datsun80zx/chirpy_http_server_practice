package api

import (
	"net/http"
	"time"

	"github.com/datsun80zx/chirpy_http_server_practice.git/internal/auth"
)

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	// Get refresh token from Authorization header
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "missing or invalid token", err)
		return
	}

	// Get user from refresh token (this also validates the token)
	user, err := h.Config.Database.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "invalid or expired refresh token", err)
		return
	}

	// Create new access token that expires in 1 hour
	accessToken, err := auth.MakeJWT(user.ID, h.Config.TokenString, time.Hour)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "problem creating token", err)
		return
	}

	RespondWithJSON(w, http.StatusOK, response{
		Token: accessToken,
	})
}
