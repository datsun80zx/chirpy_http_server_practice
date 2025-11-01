package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/datsun80zx/chirpy_http_server_practice.git/internal/auth"
	"github.com/datsun80zx/chirpy_http_server_practice.git/internal/database"
	"github.com/google/uuid"
)

type response struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Email       string    `json:"email"`
	IsChirpyRed bool      `json:"is_chirpy_red"`
}

func (h *Handler) CreateNewUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "couldn't decode parameters", err)
		return
	}
	hashedPswd, err := auth.HashPassword(params.Password)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "couldn't hash password", err)
	}

	user, err := h.Config.Database.CreateUser(r.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPswd,
	})
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Couldn't create new user", err)
		return
	}
	RespondWithJSON(w, http.StatusCreated, response{
		ID:          user.ID,
		Email:       user.Email,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		IsChirpyRed: user.IsChirpyRed.Bool,
	})
}

func (h *Handler) UpdateUserData(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "invalid jwt", err)
		return
	}
	userID, err := auth.ValidateJWT(token, h.Config.TokenString)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "not authorized user", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "couldn't decode parameters", err)
		return
	}

	hashedPswd, err := auth.HashPassword(params.Password)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "couldn't hash password", err)
	}

	newParams := database.UpdateUserPasswordAndEmailParams{
		Email:          params.Email,
		HashedPassword: hashedPswd,
		ID:             userID,
	}

	user, err := h.Config.Database.UpdateUserPasswordAndEmail(r.Context(), newParams)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Couldn't update user", err)
		return
	}
	RespondWithJSON(w, http.StatusOK, response{
		ID:          user.ID,
		Email:       user.Email,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		IsChirpyRed: user.IsChirpyRed.Bool,
	})

}

func (h *Handler) UpgradeUser(w http.ResponseWriter, r *http.Request) {
	type eventParameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "invalid key", err)
		return
	}

	if apiKey != h.Config.PolkaAPIKey {
		RespondWithError(w, http.StatusUnauthorized, "invalid key", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	event := eventParameters{}
	err = decoder.Decode(&event)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "couldn't decode parameters", err)
		return
	}
	if event.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	userID, err := uuid.Parse(event.Data.UserID)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "invalid UUID", err)
		return
	}

	_, err = h.Config.Database.GetUserByID(r.Context(), userID)
	if err != nil {
		if err == sql.ErrNoRows {
			RespondWithError(w, http.StatusNotFound, "user doesn't exist", err)
			return
		}
		RespondWithError(w, http.StatusInternalServerError, "database connection error", err)
		return
	}

	_, err = h.Config.Database.UpgradeUserAccount(r.Context(), userID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "couldn't upgrade user", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
