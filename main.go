package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

const database, collection = "url-shortener", "urls"

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

	log.Printf("Trying to retrieve the url for hash: %s", shortenedUrl)
	var url URL
	err := collection.FindOne(ctx, filter).Decode(&url)
	if err != nil {
		http.Error(w, "Url not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, url.URL, http.StatusSeeOther)
}

func PostShortenedUrl(w http.ResponseWriter, r *http.Request) {
	var payload PostShortenedUrlPayload
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	err = validateStruct(payload)
	if err != nil {
		http.Error(w, "The 'url' is missing from the payload", http.StatusBadRequest)
		return
	}

	client := NewDbClient()
	defer client.Disconnect(context.Background())

	collection := client.Database(database).Collection(collection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// TODO: replace uuid creation by a custom hash function
	hash := uuid.New()
	hashString := hash.String()

	log.Printf("Trying to insert the url: %s", payload.URL)
	_, err = collection.InsertOne(ctx, URL{
		Hash: hashString,
		URL:  payload.URL,
	})
	if err != nil {
		http.Error(w, "An error occured while saving the url", http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"hash": hashString,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("An error occured while loading the .env file")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", GetHome)
	mux.HandleFunc("GET /{shortenedUrl}", GetShortenedUrl)
	mux.HandleFunc("POST /shorten", PostShortenedUrl)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.ListenAndServe(":"+port, mux)
}
