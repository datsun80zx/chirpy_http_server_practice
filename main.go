package main

import (
	"fmt"
	"net/http"
)

func main() {

	testHandler := http.NewServeMux()

	testServ := http.Server{
		Addr:    ":8080",
		Handler: testHandler,
	}

	fmt.Println("Server starting on port 8080...")
	testServ.ListenAndServe()
}
