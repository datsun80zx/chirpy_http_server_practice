package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/datsun80zx/chirpy_http_server_practice.git/internal/auth"
	"github.com/datsun80zx/chirpy_http_server_practice.git/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) CreateNewUser(w http.ResponseWriter, r *http.Request) {
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
		respondWithError(w, http.StatusInternalServerError, "couldn't decode parameters", err)
		return
	}
	hashedPswd, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't hash password", err)
	}

	user, err := cfg.database.CreateUser(r.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPswd,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't create new user", err)
		return
	}
	respondWithJSON(w, http.StatusCreated, response{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}
