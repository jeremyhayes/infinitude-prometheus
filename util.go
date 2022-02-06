package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

func gauge(descriptor *prometheus.Desc, value float64, labelValues ...string) prometheus.Metric {
	return prometheus.MustNewConstMetric(
		descriptor,
		prometheus.GaugeValue,
		value,
		labelValues...,
	)
}

func getJson(url string, target interface{}) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return json.NewDecoder(res.Body).Decode(target)
}

func parseFloat(s string) float64 {
	f, _ := strconv.ParseFloat(s, 32)
	return f
}
