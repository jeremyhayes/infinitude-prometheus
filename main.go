package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/health", health)
	http.ListenAndServe(":8080", nil)
}

func health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "alive")
}
