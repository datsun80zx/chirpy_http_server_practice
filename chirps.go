package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/datsun80zx/chirpy_http_server_practice.git/internal/auth"
	"github.com/datsun80zx/chirpy_http_server_practice.git/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) CreateNewChirp(w http.ResponseWriter, r *http.Request) {

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid jwt", err)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.tokenString)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "not authorized user", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	chirpParams := database.CreateChirpParams{}
	err = decoder.Decode(&chirpParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode parameters", err)
		return
	}

	chirpParams.UserID = userID

	chirp, ok := validateChirp(chirpParams.Body)
	if ok {
		chirpParams.Body = chirp
		newChirp, err := cfg.database.CreateChirp(r.Context(), chirpParams)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Couldn't create chirp", err)
			return
		}
		respondWithJSON(w, http.StatusCreated, newChirp)
	} else {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

}

// this function will test the validity of a chirp. If the chirp is invalid it will return the chirp and false, otherwise it will return the string and true.
func validateChirp(chirp string) (string, bool) {
	const chirpLength = 140

	if len(chirp) > chirpLength {
		return chirp, false
	}

	return wordFilter(chirp), true

}

func wordFilter(s string) string {
	invalidWords := map[string]bool{
		"kerfuffle": true,
		"sharbert":  true,
		"fornax":    true,
	}

	wordList := strings.Split(s, " ")
	for i, word := range wordList {
		if invalidWords[strings.ToLower(word)] {
			wordList[i] = "****"
		}
	}

	return strings.Join(wordList, " ")

}

func (cfg *apiConfig) GetAllChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.database.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't create chirp", err)
		return
	}
	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) GetOneChirp(w http.ResponseWriter, r *http.Request) {
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	chirp, err := cfg.database.GetOneChirp(r.Context(), chirpID)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, "Chirp not found", nil)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Database error", err)
		return
	}
	respondWithJSON(w, http.StatusOK, chirp)
}
