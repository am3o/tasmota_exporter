package collector

import (
	"github.com/prometheus/client_golang/prometheus"
)

/*
	`device_power_total{address="` + d.address + `",label="` + d.label + `"} ` + strconv.FormatFloat(d.power.StatusSNS.ENERGY.Total, 'E', -1, 32) + "\n" +
	`device_power_today{address="` + d.address + `",label="` + d.label + `"} ` + strconv.FormatFloat(d.power.StatusSNS.ENERGY.Today, 'E', -1, 32) + "\n" +
	`device_power_yesterday{address="` + d.address + `",label="` + d.label + `"} ` + strconv.FormatFloat(d.power.StatusSNS.ENERGY.Yesterday, 'E', -1, 32) + "\n" +
	`device_power_power{address="` + d.address + `",label="` + d.label + `"} ` + strconv.Itoa(d.power.StatusSNS.ENERGY.Power) + "\n" +
	`device_power_apparent{address="` + d.address + `",label="` + d.label + `"} ` + strconv.Itoa(d.power.StatusSNS.ENERGY.ApparentPower) + "\n" +
	`device_power_reactive{address="` + d.address + `",label="` + d.label + `"} ` + strconv.Itoa(d.power.StatusSNS.ENERGY.ReactivePower) + "\n" +
	`device_power_factor{address="` + d.address + `",label="` + d.label + `"} ` + strconv.FormatFloat(d.power.StatusSNS.ENERGY.Factor, 'E', -1, 32) + "\n" +
	`device_power_voltage{address="` + d.address + `",label="` + d.label + `"} ` + strconv.Itoa(d.power.StatusSNS.ENERGY.Voltage) + "\n" +
	`device_power_current{address="` + d.address + `",label="` + d.label + `"} ` + strconv.FormatFloat(d.power.StatusSNS.ENERGY.Current, 'E', -1, 32) + "\n"
*/

const (
	Namespace    = "tasmota"
	labelSDK     = "sdk"
	labelVersion = "version"
	labelUnit    = "unit"
)

type Collector struct {
	Device              *prometheus.CounterVec
	PowerUsageTotal     *prometheus.GaugeVec
	PowerUsageToday     *prometheus.GaugeVec
	PowerUsageYesterday *prometheus.GaugeVec
	PowerUsageCurrent   *prometheus.GaugeVec
}

func New(deviceName string, ip string, deviceType string) *Collector {
	constLabels := prometheus.Labels{
		"name":        deviceName,
		"ip":          ip,
		"device_type": deviceType,
	}

	return &Collector{
		Device: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace:   Namespace,
			Name:        "device_information",
			ConstLabels: constLabels,
		}, []string{labelSDK, labelVersion}),
		PowerUsageTotal: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace:   Namespace,
			Name:        "device_power_usage_total",
			ConstLabels: constLabels,
		}, []string{labelUnit}),
		PowerUsageToday: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: Namespace,
			Name:      "device_power_usage_today",
		}, []string{labelUnit}),
		PowerUsageYesterday: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace:   Namespace,
			Name:        "device_power_usage_yesterday",
			ConstLabels: constLabels,
		}, []string{labelUnit}),
		PowerUsageCurrent: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace:   Namespace,
			Name:        "device_power_usage_current",
			ConstLabels: constLabels,
		}, []string{labelUnit}),
	}
}

func (c *Collector) Version(sdk string, version string) {
	c.Device.With(prometheus.Labels{labelSDK: sdk, labelVersion: version}).Inc()
}

func (c *Collector) Describe(descs chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(c, descs)
}

func (c *Collector) Collect(metrics chan<- prometheus.Metric) {
	c.PowerUsageTotal.Collect(metrics)
	c.PowerUsageToday.Collect(metrics)
	c.PowerUsageYesterday.Collect(metrics)
	c.PowerUsageCurrent.Collect(metrics)
	c.Device.Collect(metrics)
}

func (c *Collector) SetPowerUsage(current, today, yesterday, total float64) {
	c.PowerUsageTotal.With(prometheus.Labels{labelUnit: ""}).Set(total)
	c.PowerUsageToday.With(prometheus.Labels{labelUnit: ""}).Set(today)
	c.PowerUsageYesterday.With(prometheus.Labels{labelUnit: ""}).Set(yesterday)
	c.PowerUsageCurrent.With(prometheus.Labels{labelUnit: ""}).Set(current)
}
