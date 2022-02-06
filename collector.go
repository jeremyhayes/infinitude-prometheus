package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

// https://rsmitty.github.io/Prometheus-Exporters/

type infinitudeCollector struct {
	statusJsonUrl string
	oatMetric     *prometheus.Desc
}

func newInfinitudeCollector(baseUrl string) *infinitudeCollector {
	statusJsonUrl := fmt.Sprintf("%s/status.json", baseUrl)
	return &infinitudeCollector{
		statusJsonUrl: statusJsonUrl,
		oatMetric: prometheus.NewDesc(
			"infinitude_status_oat",
			"Outside Air Temperature",
			nil,
			nil,
		),
	}
}

func (c *infinitudeCollector) Describe(ch chan<- *prometheus.Desc) {
	// describe each metric
	ch <- c.oatMetric
}

func (c *infinitudeCollector) Collect(ch chan<- prometheus.Metric) {
	// fetch data
	statusRes := &statusResponse{}
	err := getJson(c.statusJsonUrl, statusRes)
	if err != nil {
		log.Printf("error fetching status %+v", err)
		return
	}

	// update metric values
	status := statusRes.Status[0]
	oatValue, _ := strconv.ParseFloat(status.Oat[0], 32)
	ch <- prometheus.MustNewConstMetric(c.oatMetric, prometheus.GaugeValue, oatValue)
}

func getJson(url string, target interface{}) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return json.NewDecoder(res.Body).Decode(target)
}

type statusResponse struct {
	Status []status `json:"status"`
}

type status struct {
	Version string   `json:"version"`
	Oat     []string `json:"oat"`
}
