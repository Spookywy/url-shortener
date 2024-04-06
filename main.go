package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

const database, collection = "url-shortener", "urls"

type URL struct {
	Hash string `bson:"hash"`
	URL  string `bson:"url"`
}

func GetHome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Url shortener</h1>")
}

func GetShortenedUrl(w http.ResponseWriter, r *http.Request) {
	shortenedUrl := r.PathValue("shortenedUrl")

	client := NewDbClient()
	defer client.Disconnect(context.Background())

	collection := client.Database(database).Collection(collection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := map[string]interface{}{
		"hash": shortenedUrl,
	}

	log.Printf("Trying to retrieve the url for hash %s", shortenedUrl)
	var url URL
	err := collection.FindOne(ctx, filter).Decode(&url)
	if err != nil {
		http.Error(w, "Url not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, url.URL, http.StatusSeeOther)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("An error occured while loading the .env file")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", GetHome)
	mux.HandleFunc("/{shortenedUrl}", GetShortenedUrl)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.ListenAndServe(":"+port, mux)
}
