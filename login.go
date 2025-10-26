package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/datsun80zx/chirpy_http_server_practice.git/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) Login(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password  string `json:"password"`
		Email     string `json:"email"`
		ExpiresAt int    `json:"expires_in_seconds"`
	}

	type response struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
		Token     string    `json:"token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode parameters", err)
		return
	}

	user, err := cfg.database.GetUser(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "incorrect email or password", err)
		return
	}

	isValid, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't validate info", err)
		return
	}

	if isValid {
		var expiresIn time.Duration
		if params.ExpiresAt == 0 || params.ExpiresAt > 3600 {
			expiresIn, err = time.ParseDuration("1h")
			if err != nil {
				respondWithError(w, http.StatusInternalServerError, "problem parsing expirary", err)
				return
			}
		} else {
			expiresIn, err = time.ParseDuration(fmt.Sprintf("%vs", params.ExpiresAt))
			if err != nil {
				respondWithError(w, http.StatusInternalServerError, "problem parsing expirary", err)
				return
			}
		}

		token, err := auth.MakeJWT(user.ID, cfg.tokenString, expiresIn)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "problem creating token", err)
			return
		}

		respondWithJSON(w, http.StatusOK, response{
			ID:        user.ID,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Token:     token,
		})
	} else {
		respondWithError(w, http.StatusUnauthorized, "incorrect email or password", err)
		return
	}

}
