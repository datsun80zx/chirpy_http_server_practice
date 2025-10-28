package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/datsun80zx/chirpy_http_server_practice.git/internal/auth"
	"github.com/datsun80zx/chirpy_http_server_practice.git/internal/database"
	"github.com/google/uuid"
)

func (h *Handler) CreateNewUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	type response struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
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
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}

func (h *Handler) UpdateUserData(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	type response struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
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
		RespondWithError(w, http.StatusBadRequest, "Couldn't create new user", err)
		return
	}
	RespondWithJSON(w, http.StatusOK, response{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})

}
