package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
)

const database, collection = "url-shortener", "urls"

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
	if err == mongo.ErrNoDocuments {
		http.Error(w, "Url not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "An error occured while retrieving the url", http.StatusInternalServerError)
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

	log.Printf("Verifying if the URL is already stored in the database: %s", payload.URL)
	var url URL
	err = collection.FindOne(ctx, map[string]interface{}{
		"url": payload.URL,
	}).Decode(&url)
	if err == nil {
		log.Printf("Url already stored, returning the shortened link...")
		htmlShortenedLink := fmt.Sprintf("<a href=\"/%s\">🔗 Shortened link</a>", url.Hash)

		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusFound)
		w.Write([]byte(htmlShortenedLink))
		return
	}
	if err != mongo.ErrNoDocuments {
		http.Error(w, "An error occured", http.StatusInternalServerError)
		return
	}

	log.Printf("Trying to generate a unique hash...")
	var count int64
	var hashString string
	hash := sha256.Sum256([]byte(payload.URL))
	for i := MIN_HASH_LENGTH; i <= MAX_HASH_LENGTH; i++ {
		hashString = hex.EncodeToString(hash[:])[:i]

		filter := map[string]interface{}{
			"hash": hashString,
		}
		count, err = collection.CountDocuments(ctx, filter)
		if err != nil {
			http.Error(w, "An error occured", http.StatusInternalServerError)
			return
		}
		// If we haven't used this hash before, we can use it
		if count == 0 {
			log.Printf("Unique hash found: %s", hashString)
			break
		}
		log.Printf("Attempt %d failed: %s already exists", i, hashString)
	}

	if count != 0 {
		http.Error(w, "A unique hash could not be found within the limited number of attempts", http.StatusInternalServerError)
		return
	}

	log.Printf("Trying to insert the url: %s", payload.URL)
	_, err = collection.InsertOne(ctx, URL{
		Hash: hashString,
		URL:  payload.URL,
	})
	if err != nil {
		http.Error(w, "An error occured while saving the url", http.StatusInternalServerError)
		return
	}

	htmlShortenedLink := fmt.Sprintf("<a href=\"/%s\">🔗 Shortened link</a>", hashString)

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(htmlShortenedLink))
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("An error occured while loading the .env file")
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("public")))
	mux.HandleFunc("GET /{shortenedUrl}", GetShortenedUrl)
	mux.HandleFunc("POST /shorten", PostShortenedUrl)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.ListenAndServe(":"+port, mux)
}
