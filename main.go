package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func GetHome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Url shortener</h1>")
}

func GetShortenedUrl(w http.ResponseWriter, r *http.Request) {
	shortenedUrl := r.PathValue("shortenedUrl")
	fmt.Fprintf(w, "<h1>You tried to access the following url: %s</h1>", shortenedUrl)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("An error occured while loading the .env file")
	}

	client := NewDbClient()
	// Remove the following lines (test purposes)
	colection := client.Database("url_shortener").Collection("urls")
	fmt.Println(colection)

	mux := http.NewServeMux()
	mux.HandleFunc("/", GetHome)
	mux.HandleFunc("/{shortenedUrl}", GetShortenedUrl)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.ListenAndServe(":"+port, mux)
}
