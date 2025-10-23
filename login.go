package main

import (
	"encoding/json"
	"net/http"

	"github.com/datsun80zx/chirpy_http_server_practice.git/internal/auth"
	"github.com/datsun80zx/chirpy_http_server_practice.git/internal/database"
)

func (cfg *apiConfig) Login(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	type response struct {
		database.User
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
		respondWithJSON(w, http.StatusOK, response{
			User: database.User{
				ID:        user.ID,
				Email:     user.Email,
				CreatedAt: user.CreatedAt,
				UpdatedAt: user.UpdatedAt,
			},
		})
	} else {
		respondWithError(w, http.StatusUnauthorized, "incorrect email or password", err)
		return
	}

}
