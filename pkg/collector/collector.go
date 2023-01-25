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

	labelDeviceName = "deviceName"
	labelIP         = "ip"

	labelDeviceType = "deviceType"
)

type Metadata struct {
	IP     string
	Device Device
}

type Device struct {
	Name string
	Type string
}
type Collector struct {
	Device              *prometheus.CounterVec
	PowerUsageTotal     *prometheus.GaugeVec
	PowerUsageToday     *prometheus.GaugeVec
	PowerUsageYesterday *prometheus.GaugeVec
	PowerUsageCurrent   *prometheus.GaugeVec
}

func New() *Collector {
	return &Collector{
		Device: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: Namespace,
			Name:      "device_information",
		}, []string{labelSDK, labelVersion, labelDeviceName, labelIP, labelDeviceType}),
		PowerUsageTotal: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: Namespace,
			Name:      "device_power_usage_total",
		}, []string{labelUnit, labelDeviceName, labelIP, labelDeviceType}),
		PowerUsageToday: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: Namespace,
			Name:      "device_power_usage_today",
		}, []string{labelUnit, labelDeviceName, labelIP, labelDeviceType}),
		PowerUsageYesterday: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: Namespace,
			Name:      "device_power_usage_yesterday",
		}, []string{labelUnit, labelDeviceName, labelIP, labelDeviceType}),
		PowerUsageCurrent: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: Namespace,
			Name:      "device_power_usage_current",
		}, []string{labelUnit, labelDeviceName, labelIP, labelDeviceType}),
	}
}

func (c *Collector) Version(sdk string, version string, metadata Metadata) {
	c.Device.With(
		prometheus.Labels{
			labelSDK:        sdk,
			labelVersion:    version,
			labelDeviceName: metadata.Device.Name,
			labelIP:         metadata.IP,
			labelDeviceType: metadata.Device.Type,
		}).Inc()
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

func (c *Collector) SetPowerUsage(current, today, yesterday, total float64, metadata Metadata) {
	c.setPowerUsageTotal(total, metadata)
	c.setPowerUsageToday(today, metadata)
	c.setPowerUsageYesterday(yesterday, metadata)
	c.setPowerUsageCurrent(current, metadata)
}

func (c *Collector) setPowerUsageTotal(value float64, metadata Metadata) {
	c.PowerUsageTotal.With(prometheus.Labels{
		labelUnit:       "",
		labelIP:         metadata.IP,
		labelDeviceName: metadata.Device.Name,
		labelDeviceType: metadata.Device.Type,
	}).Set(value)
}
func (c *Collector) setPowerUsageToday(value float64, metadata Metadata) {
	c.PowerUsageToday.With(prometheus.Labels{
		labelUnit:       "",
		labelIP:         metadata.IP,
		labelDeviceName: metadata.Device.Name,
		labelDeviceType: metadata.Device.Type,
	}).Set(value)
}

func (c *Collector) setPowerUsageYesterday(value float64, metadata Metadata) {
	c.PowerUsageYesterday.With(prometheus.Labels{
		labelUnit:       "",
		labelIP:         metadata.IP,
		labelDeviceName: metadata.Device.Name,
		labelDeviceType: metadata.Device.Type,
	}).Set(value)
}

func (c *Collector) setPowerUsageCurrent(value float64, metadata Metadata) {
	c.PowerUsageCurrent.With(prometheus.Labels{
		labelUnit:       "",
		labelIP:         metadata.IP,
		labelDeviceName: metadata.Device.Name,
		labelDeviceType: metadata.Device.Type,
	}).Set(value)
}
