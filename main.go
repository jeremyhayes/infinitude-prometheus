package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// disable default process and go metrics collectors, too noisy
	prometheus.Unregister(collectors.NewGoCollector())
	prometheus.Unregister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{
		PidFn: func() (int, error) {
			return os.Getpid(), nil
		},
		Namespace: "",
	}))

	infinitudeBaseUrl := os.Getenv("INFINITUDE_BASE_URL")
	infinitudeCollector := newInfinitudeCollector(infinitudeBaseUrl)
	prometheus.Register(infinitudeCollector)

	http.HandleFunc("/health", health)
	http.Handle("/metrics", promhttp.Handler())

	log.Println("Starting http server on :8080...")
	log.Printf("Monitoring Infinitude instance at %s\n", infinitudeBaseUrl)
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
