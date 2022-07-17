package collectors

import (
	"fmt"
	"log"

	"github.com/prometheus/client_golang/prometheus"
)

// https://rsmitty.github.io/Prometheus-Exporters/

type energyCollector struct {
	energyJsonUrl string
	//
	usageEmergencyHeat *prometheus.Desc
	usageFanGas        *prometheus.Desc
	usageReheat        *prometheus.Desc
	usageFan           *prometheus.Desc
	usageCooling       *prometheus.Desc
	usageLoopPump      *prometheus.Desc
	usageHeatPumpHeat  *prometheus.Desc
	usageGas           *prometheus.Desc
	//
	costEmergencyHeat *prometheus.Desc
	costFanGas        *prometheus.Desc
	costReheat        *prometheus.Desc
	costFan           *prometheus.Desc
	costCooling       *prometheus.Desc
	costLoopPump      *prometheus.Desc
	costHeatPumpHeat  *prometheus.Desc
	costGas           *prometheus.Desc
}

func NewEnergyCollector(baseUrl string) *energyCollector {
	energyJsonUrl := fmt.Sprintf("%s/energy.json", baseUrl)
	return &energyCollector{
		energyJsonUrl: energyJsonUrl,
		// usage metrics
		usageEmergencyHeat: prometheus.NewDesc(
			"infinitude_energy_usage_eheat",
			"",
			[]string{"period_id"},
			nil,
		),
		usageFanGas: prometheus.NewDesc(
			"infinitude_energy_usage_fangas",
			"",
			[]string{"period_id"},
			nil,
		),
		usageReheat: prometheus.NewDesc(
			"infinitude_energy_usage_reheat",
			"",
			[]string{"period_id"},
			nil,
		),
		usageFan: prometheus.NewDesc(
			"infinitude_energy_usage_fan",
			"",
			[]string{"period_id"},
			nil,
		),
		usageCooling: prometheus.NewDesc(
			"infinitude_energy_usage_cooling",
			"",
			[]string{"period_id"},
			nil,
		),
		usageLoopPump: prometheus.NewDesc(
			"infinitude_energy_usage_looppump",
			"",
			[]string{"period_id"},
			nil,
		),
		usageHeatPumpHeat: prometheus.NewDesc(
			"infinitude_energy_usage_hpheat",
			"",
			[]string{"period_id"},
			nil,
		),
		usageGas: prometheus.NewDesc(
			"infinitude_energy_usage_gas",
			"",
			[]string{"period_id"},
			nil,
		),
		// cost metrics
		costEmergencyHeat: prometheus.NewDesc(
			"infinitude_energy_cost_eheat",
			"",
			[]string{"period_id"},
			nil,
		),
		costFanGas: prometheus.NewDesc(
			"infinitude_energy_cost_fangas",
			"",
			[]string{"period_id"},
			nil,
		),
		costReheat: prometheus.NewDesc(
			"infinitude_energy_cost_reheat",
			"",
			[]string{"period_id"},
			nil,
		),
		costFan: prometheus.NewDesc(
			"infinitude_energy_cost_fan",
			"",
			[]string{"period_id"},
			nil,
		),
		costCooling: prometheus.NewDesc(
			"infinitude_energy_cost_cooling",
			"",
			[]string{"period_id"},
			nil,
		),
		costLoopPump: prometheus.NewDesc(
			"infinitude_energy_cost_looppump",
			"",
			[]string{"period_id"},
			nil,
		),
		costHeatPumpHeat: prometheus.NewDesc(
			"infinitude_energy_cost_hpheat",
			"",
			[]string{"period_id"},
			nil,
		),
		costGas: prometheus.NewDesc(
			"infinitude_energy_cost_gas",
			"",
			[]string{"period_id"},
			nil,
		),
	}
}

func (c *energyCollector) Describe(ch chan<- *prometheus.Desc) {
	// describe each metric

	// usage metrics
	ch <- c.usageEmergencyHeat
	ch <- c.usageFanGas
	ch <- c.usageReheat
	ch <- c.usageFan
	ch <- c.usageCooling
	ch <- c.usageLoopPump
	ch <- c.usageHeatPumpHeat
	ch <- c.usageGas

	// cost metrics
	ch <- c.costEmergencyHeat
	ch <- c.costFanGas
	ch <- c.costReheat
	ch <- c.costFan
	ch <- c.costCooling
	ch <- c.costLoopPump
	ch <- c.costHeatPumpHeat
	ch <- c.costGas
}

func (c *energyCollector) Collect(ch chan<- prometheus.Metric) {
	// fetch data
	resp := &EnergyResponse{}
	err := getJson(c.energyJsonUrl, resp)
	if err != nil {
		log.Printf("error fetching energy %+v", err)
		return
	}

	// check empty response
	if len(resp.Energy) == 0 {
		log.Print("no energy data returned")
		return
	}

	// update metric values
	energy := resp.Energy[0]

	for _, usage := range energy.Usage {
		for _, period := range usage.Periods {
			id := period.Id
			ch <- gauge(c.usageEmergencyHeat, parseFloat(period.EmergencyHeat[0]), id)
			ch <- gauge(c.usageFanGas, parseFloat(period.FanGas[0]), id)
			ch <- gauge(c.usageReheat, parseFloat(period.Reheat[0]), id)
			ch <- gauge(c.usageFan, parseFloat(period.Fan[0]), id)
			ch <- gauge(c.usageCooling, parseFloat(period.Cooling[0]), id)
			ch <- gauge(c.usageLoopPump, parseFloat(period.LoopPump[0]), id)
			ch <- gauge(c.usageHeatPumpHeat, parseFloat(period.HeatPumpHeat[0]), id)
			ch <- gauge(c.usageGas, parseFloat(period.Gas[0]), id)
		}
	}
	for _, cost := range energy.Cost {
		for _, period := range cost.Periods {
			id := period.Id
			ch <- gauge(c.costEmergencyHeat, parseFloat(period.EmergencyHeat[0]), id)
			ch <- gauge(c.costFanGas, parseFloat(period.FanGas[0]), id)
			ch <- gauge(c.costReheat, parseFloat(period.Reheat[0]), id)
			ch <- gauge(c.costFan, parseFloat(period.Fan[0]), id)
			ch <- gauge(c.costCooling, parseFloat(period.Cooling[0]), id)
			ch <- gauge(c.costLoopPump, parseFloat(period.LoopPump[0]), id)
			ch <- gauge(c.costHeatPumpHeat, parseFloat(period.HeatPumpHeat[0]), id)
			ch <- gauge(c.costGas, parseFloat(period.Gas[0]), id)
		}
	}
}

type EnergyResponse struct {
	Energy []Energy `json:"energy"`
}

type Energy struct {
	Version string   `json:"version"`
	hspf    []string `json:"hspf"` // Heating Seasonal Performance Factor
	seer    []string `json:"seer"` // Seasonal Energy Efficiency Ratio
	Usage   []Usage  `json:"usage"`
	Cost    []Cost   `json:"cost"`

	// cooling
	// fan
	// reheat
	// looppump
	// hpheat
	// gas
	// eheat
	// fangas
}

type Usage struct {
	Periods []UsagePeriod `json:"period"`
}

type UsagePeriod struct {
	Id            string   `json:"id"`
	Cooling       []string `json:"cooling"`
	EmergencyHeat []string `json:"eheat"`
	Fan           []string `json:"fan"`
	FanGas        []string `json:"fangas"`
	Gas           []string `json:"gas"`
	HeatPumpHeat  []string `json:"hpheat"`
	LoopPump      []string `json:"looppump"`
	Reheat        []string `json:"reheat"`
}

type Cost struct {
	Periods []CostPeriod `json:"period"`
}

type CostPeriod struct {
	Id            string   `json:"id"`
	Cooling       []string `json:"cooling"`
	EmergencyHeat []string `json:"eheat"`
	Fan           []string `json:"fan"`
	FanGas        []string `json:"fangas"`
	Gas           []string `json:"gas"`
	HeatPumpHeat  []string `json:"hpheat"`
	LoopPump      []string `json:"looppump"`
	Reheat        []string `json:"reheat"`
}
