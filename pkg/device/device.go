package device

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	DEVICE_INFORMATION = "status"
	DEVICE_VERSION     = "status+2"
	DEVICE_NETWORK     = "status+5"
	DEVICE_ENERGY      = "status+8"
)

type PowerDevice struct {
	url url.URL
}

func New(address string, username string, password string) PowerDevice {
	return PowerDevice{url: url.URL{
		Scheme: "http",
		Host:   address,
	}}
}

func (p PowerDevice) Version(ctx context.Context) (Version, error) {
	p.url.Path = "/cm"

	query := p.url.Query()
	query.Set("cmnd", DEVICE_VERSION)
	p.url.RawQuery = query.Encode()

	req, err := http.NewRequest(http.MethodGet, p.url.String(), nil)
	if err != nil {
		return Version{}, err
	}
	req = req.WithContext(ctx)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Version{}, err
	}

	if err != nil || resp.StatusCode != http.StatusOK {
		return Version{}, fmt.Errorf("could not fetch device: %w", err)
	}
	defer resp.Body.Close()

	var version Version
	if err := json.NewDecoder(resp.Body).Decode(&version); err != nil {
		return Version{}, err
	}

	return version, nil
}
func (p PowerDevice) Network(ctx context.Context) (Network, error) {
	p.url.Path = "/cm"

	query := p.url.Query()
	query.Set("cmnd", DEVICE_NETWORK)
	p.url.RawQuery = query.Encode()

	req, err := http.NewRequest(http.MethodGet, p.url.String(), nil)
	if err != nil {
		return Network{}, err
	}
	req = req.WithContext(ctx)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Network{}, err
	}

	if err != nil || resp.StatusCode != http.StatusOK {
		return Network{}, fmt.Errorf("could not fetch device: %w", err)
	}
	defer resp.Body.Close()

	var network Network
	if err := json.NewDecoder(resp.Body).Decode(&network); err != nil {
		return Network{}, err
	}

	return network, nil
}

func (p PowerDevice) DeviceInformation(ctx context.Context) (Device, error) {
	p.url.Path = "/cm"

	query := p.url.Query()
	query.Set("cmnd", DEVICE_INFORMATION)
	p.url.RawQuery = query.Encode()

	req, err := http.NewRequest(http.MethodGet, p.url.String(), nil)
	if err != nil {
		return Device{}, err
	}
	req = req.WithContext(ctx)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Device{}, err
	}

	if err != nil || resp.StatusCode != http.StatusOK {
		return Device{}, fmt.Errorf("could not fetch device: %w", err)
	}
	defer resp.Body.Close()

	var device Device
	if err := json.NewDecoder(resp.Body).Decode(&device); err != nil {
		return Device{}, err
	}

	return device, nil
}

func (p PowerDevice) Status(ctx context.Context) (PowerStatus, error) {
	p.url.Path = "/cm"

	query := p.url.Query()
	query.Set("cmnd", DEVICE_ENERGY)
	p.url.RawQuery = query.Encode()

	req, err := http.NewRequest(http.MethodGet, p.url.String(), nil)
	if err != nil {
		return PowerStatus{}, err
	}
	req = req.WithContext(ctx)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return PowerStatus{}, err
	}

	if err != nil || resp.StatusCode != http.StatusOK {
		return PowerStatus{}, fmt.Errorf("could not fetch device: %w", err)
	}
	defer resp.Body.Close()

	var status PowerStatus
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return PowerStatus{}, err
	}

	return status, nil
}
