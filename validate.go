package main

import (
	"encoding/json"
	"net/http"
)

// check if json body is 140char or less
func validateChirp(w http.ResponseWriter, r *http.Request) {

	type chirpBody struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	chirp := chirpBody{}
	err := decoder.Decode(&chirp)
	if err != nil {
		// err handling logic
	}
}

//if under 140char respond with 200 ok

// else respond 400 status w/ json body explaining error
