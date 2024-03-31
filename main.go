package main

import (
	"fmt"
	"net/http"
	"os"
)

func GetShortenedUrl(w http.ResponseWriter, r *http.Request) {
	shortenedUrl := r.PathValue("shortenedUrl")
	fmt.Fprintf(w, "<h1>You tried to access the following url: %s</h1>", shortenedUrl)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/{shortenedUrl}", GetShortenedUrl)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.ListenAndServe(":"+port, mux)
}
