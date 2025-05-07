package main

import (
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	// Serve static files (like index.html) from the current directory (.)
	fileServer := http.FileServer(http.Dir("."))

	// Register the root path ("/") with the file server
	mux.Handle("/", fileServer)

	fileServerAssets := http.FileServer(http.Dir("./assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", fileServerAssets))

	// Create and start the server
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	server.ListenAndServe()
}
