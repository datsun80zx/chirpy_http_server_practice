package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// check if json body is 140char or less
func validateChirp(w http.ResponseWriter, r *http.Request) {

	type chirpBody struct {
		Body string `json:"body"`
	}

	type badRequest struct {
		Error string `json:"error"`
	}

	type successfulRequest struct {
		Valid bool `json:"valid"`
	}

	decoder := json.NewDecoder(r.Body)
	chirp := chirpBody{}
	err := decoder.Decode(&chirp)
	if err != nil {
		badReq := badRequest{
			Error: fmt.Sprintf("%v", err),
		}
		dat, err := json.Marshal(badReq)
		if err != nil {
			log.Printf("Error marshalling json: %s", err)
			w.WriteHeader(500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write(dat)
		return
	}

	if len(chirp.Body) > 140 {
		log.Println("Chirp is too long! Length:", len(chirp.Body))
		badReq := badRequest{
			Error: "Chirp is too long",
		}
		dat, err := json.Marshal(badReq)
		if err != nil {
			log.Printf("Error marshalling json: %s", err)
			w.WriteHeader(500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write(dat)
		return
	} else {
		succResp := successfulRequest{
			Valid: true,
		}
		dat, err := json.Marshal(succResp)
		if err != nil {
			log.Printf("Error marshalling json: %s", err)
			w.WriteHeader(500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(dat)
	}

}
