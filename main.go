package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/am3o/tasmota_exporter/pkg/collector"
	"github.com/am3o/tasmota_exporter/pkg/config"
	"github.com/am3o/tasmota_exporter/pkg/device"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

var configFile, _ = os.LookupEnv("CONFIG_FILE")

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	configuration, err := config.New(configFile)
	if err != nil {
		logger.Error("could not read configuration", zap.Error(err))
		os.Exit(1)
		return
	}

	metrics := collector.New()

	ctx := context.Background()
	devices := make([]device.PowerDevice, 0)
	for _, value := range configuration.Devices {
		sensor := device.New(value.IP, value.Username, value.Password, value.Name, value.Type)
		version, err := sensor.Version(ctx)
		if err != nil {
			logger.With(
				zap.String("ip", value.IP),
				zap.String("device", value.Name),
				zap.String("type", value.Type),
			).Error("could not detect version of device", zap.Error(err))
			os.Exit(1)
			return
		}

		metrics.Version(version.Status.SDK, version.Status.Version, collector.Metadata{
			IP: value.IP,
			Device: collector.Device{
				Name: value.Name,
				Type: value.Type,
			},
		})
		devices = append(devices, sensor)
	}
	prometheus.MustRegister(metrics)

	logger.Info("service is initialized and starts working")
	go func(ctx context.Context) {
		ticker := time.NewTicker(time.Second * 15)
		for ; ; <-ticker.C {
			for _, value := range devices {
				go func(powerDevice device.PowerDevice) {
					info, err := powerDevice.Status(ctx)
					if err != nil {
						logger.With(
							zap.String("ip", powerDevice.Information.IP),
							zap.String("device", powerDevice.Information.Name),
							zap.String("type", powerDevice.Information.Type),
						).Error("could not fetch information from the device", zap.Error(err))
						return
					}

					metrics.SetPowerUsage(
						info.Status.Energy.Current,
						info.Status.Energy.Today,
						info.Status.Energy.Yesterday,
						info.Status.Energy.Total,
						collector.Metadata{
							IP: powerDevice.Information.IP,
							Device: collector.Device{
								Name: powerDevice.Information.Name,
								Type: powerDevice.Information.Type,
							},
						},
					)

					logger.With(
						zap.String("device", powerDevice.Information.Name),
						zap.String("ip", powerDevice.Information.IP),
						zap.String("type", powerDevice.Information.Type),
					).Info("update all values from device")
				}(value)
			}
		}
	}(context.Background())

	http.Handle("/internal/metrics", promhttp.Handler())
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
