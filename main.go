package main

import (
	"log"
	"net/http"
)

func main() {
	const port = "8080"
	const filepathRoot = "."

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(filepathRoot)))

	testServ := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Server starting...\n\n")
	log.Printf("Now serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(testServ.ListenAndServe())

}
