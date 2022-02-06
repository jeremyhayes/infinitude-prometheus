package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

// https://rsmitty.github.io/Prometheus-Exporters/

type infinitudeCollector struct {
	oatMetric *prometheus.Desc
}

func newInfinitudeCollector() *infinitudeCollector {
	return &infinitudeCollector{
		oatMetric: prometheus.NewDesc(
			"infinitude_oat",
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
	// update metric values
	gaugeValue := float64(42)
	ch <- prometheus.MustNewConstMetric(c.oatMetric, prometheus.GaugeValue, gaugeValue)
}
