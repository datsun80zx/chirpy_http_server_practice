package main

import (
	"log"
	"net/http"
)

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func main() {
	const port = "8080"
	const filepathRoot = "."

	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))

	mux.HandleFunc("/healthz", healthzHandler)

	testServ := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Server starting...\n\n")
	log.Printf("Now serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(testServ.ListenAndServe())

}
