package main

type URL struct {
	Hash string `bson:"hash"`
	URL  string `bson:"url"`
}

type PostShortenedUrlPayload struct {
	URL string `json:"url" validate:"required"`
}
