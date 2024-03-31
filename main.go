package main

import (
	"fmt"
	"net/http"
	"os"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>Hello from Go!</h1>")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", Handler)

	port := os.Getenv("PORT")
	http.ListenAndServe("0.0.0.0:"+port, mux)
}
