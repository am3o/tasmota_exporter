package main

import (
	"context"
	"github.com/am3o/tasmota_exporter/pkg/collector"
	"github.com/am3o/tasmota_exporter/pkg/device"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"net/http"
	"os"
	"time"
)

var deviceType, _ = os.LookupEnv("DEVICE_NAME")
var ipAddress, _ = os.LookupEnv("DEVICE_IP_ADDRESS")
var username, _ = os.LookupEnv("DEVICE_USERNAME")
var password, _ = os.LookupEnv("DEVICE_PASSWORD")

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	tasmota := device.New(ipAddress, username, password)
	version, err := tasmota.Version(context.Background())
	if err != nil {
		panic(err)
	}

	info, err := tasmota.Network(context.Background())
	if err != nil {
		panic(err)
	}

	metrics := collector.New(info.Status.Hostname, info.Status.Address, deviceType)
	metrics.Version(version.Status.SDK, version.Status.Version)
	prometheus.MustRegister(metrics)

	go func(ctx context.Context) {
		ticker := time.NewTicker(time.Second * 15)
		for ; ; <-ticker.C {
			info, err := tasmota.Status(ctx)
			if err != nil {
				logger.Error("could not fetch information from the device", zap.Error(err), zap.String("ip-address", ipAddress))
				continue
			}

			metrics.SetPowerUsage(
				info.Status.Energy.Current,
				info.Status.Energy.Today,
				info.Status.Energy.Yesterday,
				info.Status.Energy.Total,
			)
		}
	}(context.Background())

	http.Handle("/internal/metrics", promhttp.Handler())
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
