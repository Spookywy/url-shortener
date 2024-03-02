package main

import (
	"fmt"
	"log"
	"net/http"
)

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello there!")
}

func main() {
	log.Default().Println("Starting server")

	http.HandleFunc("/", hello)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
