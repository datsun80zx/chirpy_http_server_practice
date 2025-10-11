package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

// check if json body is 140char or less
func validateChirp(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Body string `json:"body"`
	}

	type respone struct {
		CleanedBody string `json:"cleaned_body"`
	}
	// type returnVals struct {
	// 	Valid bool `json:"valid"`
	// }

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode parameters", err)
		return
	}

	const chirpLength = 140
	if len(params.Body) > chirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	cleanedWords := wordFilter(params.Body)

	respondWithJSON(w, http.StatusOK, respone{
		CleanedBody: cleanedWords,
	})

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
