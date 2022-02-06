package main

import (
	"fmt"
	"log"

	"github.com/prometheus/client_golang/prometheus"
)

// https://rsmitty.github.io/Prometheus-Exporters/

type statusCollector struct {
	statusJsonUrl     string
	oatMetric         *prometheus.Desc
	filterLevelMetric *prometheus.Desc
	rtMetric          *prometheus.Desc
	rhMetric          *prometheus.Desc
	htspMetric        *prometheus.Desc
	clspMetric        *prometheus.Desc
}

func newStatusCollector(baseUrl string) *statusCollector {
	statusJsonUrl := fmt.Sprintf("%s/status.json", baseUrl)
	return &statusCollector{
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

func (c *statusCollector) Describe(ch chan<- *prometheus.Desc) {
	// describe each metric
	ch <- c.oatMetric
	ch <- c.filterLevelMetric
	ch <- c.rtMetric
	ch <- c.rhMetric
	ch <- c.htspMetric
	ch <- c.clspMetric
}

func (c *statusCollector) Collect(ch chan<- prometheus.Metric) {
	// fetch data
	statusRes := &StatusResponse{}
	err := getJson(c.statusJsonUrl, statusRes)
	if err != nil {
		log.Printf("error fetching status %+v", err)
		return
	}

	// update metric values
	status := statusRes.Status[0]
	ch <- gauge(c.oatMetric, parseFloat(status.Oat[0]))
	ch <- gauge(c.filterLevelMetric, parseFloat(status.FiltrLvl[0]))

	for _, zones := range status.Zones {
		zone := zones.Zone[0]
		id := zone.Id
		enabled := zone.Enabled[0] == "on"
		if !enabled {
			continue
		}

		ch <- gauge(c.rhMetric, parseFloat(zone.Rh[0]), id)

		// rt is an empty object when zone not enabled
		if rt, err := zone.Rt[0].(string); err {
			ch <- gauge(c.rtMetric, parseFloat(rt), id)
		} else {
			ch <- gauge(c.rtMetric, 0, id)
		}

		ch <- gauge(c.htspMetric, parseFloat(zone.Htsp[0]), id)
		ch <- gauge(c.clspMetric, parseFloat(zone.Clsp[0]), id)
	}

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
