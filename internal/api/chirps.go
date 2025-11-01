package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/datsun80zx/chirpy_http_server_practice.git/internal/auth"
	"github.com/datsun80zx/chirpy_http_server_practice.git/internal/database"
	"github.com/google/uuid"
)

func (h *Handler) CreateNewChirp(w http.ResponseWriter, r *http.Request) {

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
	chirpParams := database.CreateChirpParams{}
	err = decoder.Decode(&chirpParams)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "couldn't decode parameters", err)
		return
	}

	chirpParams.UserID = userID

	chirp, ok := validateChirp(chirpParams.Body)
	if ok {
		chirpParams.Body = chirp
		newChirp, err := h.Config.Database.CreateChirp(r.Context(), chirpParams)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, "Couldn't create chirp", err)
			return
		}
		RespondWithJSON(w, http.StatusCreated, newChirp)
	} else {
		RespondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
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

func (h *Handler) GetAllChirps(w http.ResponseWriter, r *http.Request) {
	author := r.URL.Query().Get("author_id")
	if author != "" {
		authorID, err := uuid.Parse(author)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, "invalid author id", err)
			return
		}
		authorChirps, err := h.Config.Database.GetChirpsByUser(r.Context(), authorID)
		if err != nil {
			if err == sql.ErrNoRows {
				RespondWithError(w, http.StatusNotFound, "no chirps available", err)
				return
			}
			RespondWithError(w, http.StatusNotFound, "no chirps available", err)
			return
		}
		RespondWithJSON(w, http.StatusOK, authorChirps)
		return
	}
	chirps, err := h.Config.Database.GetChirps(r.Context())
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Couldn't create chirp", err)
		return
	}
	RespondWithJSON(w, http.StatusOK, chirps)
}

func (h *Handler) GetOneChirp(w http.ResponseWriter, r *http.Request) {
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	chirp, err := h.Config.Database.GetOneChirp(r.Context(), chirpID)
	if err != nil {
		if err == sql.ErrNoRows {
			RespondWithError(w, http.StatusNotFound, "Chirp not found", nil)
			return
		}
		RespondWithError(w, http.StatusInternalServerError, "Database error", err)
		return
	}
	RespondWithJSON(w, http.StatusOK, chirp)
}

func (h *Handler) DeleteOneChirp(w http.ResponseWriter, r *http.Request) {
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

	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	chirp, err := h.Config.Database.GetOneChirp(r.Context(), chirpID)
	if err != nil {
		if err == sql.ErrNoRows {
			RespondWithError(w, http.StatusNotFound, "Chirp not found", nil)
			return
		}
		RespondWithError(w, http.StatusInternalServerError, "Database error", err)
		return
	}

	if chirp.UserID == userID {
		h.Config.Database.DeleteOneChirp(r.Context(), database.DeleteOneChirpParams{
			ID:     chirpID,
			UserID: userID,
		})
		w.WriteHeader(http.StatusNoContent)
	}
	RespondWithError(w, http.StatusForbidden, "not authorized", err)
}
