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
	statusJsonUrl     string
	oatMetric         *prometheus.Desc
	filterLevelMetric *prometheus.Desc
	rtMetric          *prometheus.Desc
	rhMetric          *prometheus.Desc
	htspMetric        *prometheus.Desc
	clspMetric        *prometheus.Desc
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
		filterLevelMetric: prometheus.NewDesc(
			"infinitude_status_filtrlvl",
			"Filter Level",
			nil, nil,
		),
		rtMetric: prometheus.NewDesc(
			"infinitude_status_rt",
			"Relative Temperature",
			[]string{"zone_id"},
			nil,
		),
		rhMetric: prometheus.NewDesc(
			"infinitude_status_rh",
			"Relative Humidity",
			[]string{"zone_id"},
			nil,
		),
		htspMetric: prometheus.NewDesc(
			"infinitude_status_htsp",
			"Heat Setpoint",
			[]string{"zone_id"},
			nil,
		),
		clspMetric: prometheus.NewDesc(
			"infinitude_status_clsp",
			"Cool Setpoint",
			[]string{"zone_id"},
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
	statusRes := &StatusResponse{}
	err := getJson(c.statusJsonUrl, statusRes)
	if err != nil {
		log.Printf("error fetching status %+v", err)
		return
	}

	// update metric values
	status := statusRes.Status[0]

	oatValue, _ := strconv.ParseFloat(status.Oat[0], 32)
	ch <- gauge(c.oatMetric, oatValue)

	filterLevelValue, _ := strconv.ParseFloat(status.FiltrLvl[0], 32)
	ch <- gauge(c.filterLevelMetric, filterLevelValue)

	for _, zones := range status.Zones {
		zone := zones.Zone[0]
		id := zone.Id
		enabled := zone.Enabled[0] == "on"
		if !enabled {
			continue
		}

		rhValue, _ := strconv.ParseFloat(zone.Rh[0], 32)
		ch <- gauge(c.rhMetric, rhValue, id)

		// rt is an empty object when zone not enabled
		var rtValue float64
		if rt, err := zone.Rt[0].(string); err {
			rtValue, _ = strconv.ParseFloat(rt, 32)
		} else {
			rtValue = 0
		}
		ch <- gauge(c.rtMetric, rtValue, id)

		htspValue, _ := strconv.ParseFloat(zone.Htsp[0], 32)
		ch <- gauge(c.htspMetric, htspValue, id)

		clspValue, _ := strconv.ParseFloat(zone.Clsp[0], 32)
		ch <- gauge(c.clspMetric, clspValue, id)
	}

}

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

type StatusResponse struct {
	Status []Status `json:"status"`
}

type Status struct {
	Version      string   `json:"version"`
	Oat          []string `json:"oat"`
	CfgType      []string `json:"cfgtype"`
	UvLvl        []string `json:"uvlvl"`
	Humid        []string `json:"humid"`
	OprStsMsg    []string `json:"oprstsmsg"`
	LocalTime    []string `json:"localTime"`
	CfgEm        []string `json:"cfgem"`
	HumLvl       []string `json:"humlvl`
	Zones        []Zones  `json:"zones"`
	VacatRunning []string `json:"vacatrunning"`
	Mode         []string `json:"mode"`
	FiltrLvl     []string `json:"filtrlvl"`
	VentLvl      []string `json:"ventlvl`
}

type Zones struct {
	Zone []Zone `json:"zone"`
}

type Zone struct {
	Id               string        `json:"id"`               //
	Enabled          []string      `json:"enabled"`          // on/off
	Name             []string      `json:"name"`             //
	CurrentActivity  []string      `json:"currentActivity"`  // home
	Htsp             []string      `json:"htsp"`             // Heat Setpoint
	Clsp             []string      `json:"clsp"`             // Cool Setpoint
	Rt               []interface{} `json:"rt"`               // Relative Temperature
	Rh               []string      `json:"rh"`               // Relative Humidity
	Fan              []string      `json:"fan"`              // low
	Hold             []string      `json:"hold"`             // on/off
	DamperPosition   []string      `json:"damperposition`    // 15
	ZoneConditioning []string      `json:"zoneconditioning"` // active_heat
	Otmr             []interface{} `json:"otmr"`             //
}
