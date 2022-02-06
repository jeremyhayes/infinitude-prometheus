package main

import (
	"encoding/json"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	http.HandleFunc("/health", health)
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":8080", nil)
}

func health(w http.ResponseWriter, r *http.Request) {
	status := HealthResponse{
		Status: "healthy",
	}
	json.NewEncoder(w).Encode(status)
}

type HealthResponse struct {
	Status string `json:"status"`
}
