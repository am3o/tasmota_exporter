package device

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	DEVICE_INFORMATION = "status"
	DEVICE_VERSION     = "status+2"
	DEVICE_NETWORK     = "status+5"
	DEVICE_ENERGY      = "status+8"
)

type Metadata struct {
	IP   string
	Name string
	Type string
}

type PowerDevice struct {
	Information Metadata
	url         url.URL
}

func New(address string, username string, password string, name string, Type string) PowerDevice {
	return PowerDevice{
		Information: Metadata{
			IP:   address,
			Name: name,
			Type: Type,
		},
		url: url.URL{
			Scheme: "http",
			Host:   address,
		},
	}
}

func (p PowerDevice) executeRequest(ctx context.Context, command string) (io.ReadCloser, error) {
	p.url.Path = "/cm"

	query := p.url.Query()
	query.Set("cmnd", command)
	p.url.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.url.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("could not fetch device: %w", err)
	}
	return io.NopCloser(resp.Body), nil
}

func (p PowerDevice) Version(ctx context.Context) (Version, error) {
	content, err := p.executeRequest(ctx, DEVICE_VERSION)
	if err != nil {
		return Version{}, err
	}
	defer content.Close()

	var version Version
	if err := json.NewDecoder(content).Decode(&version); err != nil {
		return Version{}, err
	}

	return version, nil
}

func (p PowerDevice) Network(ctx context.Context) (Network, error) {
	content, err := p.executeRequest(ctx, DEVICE_NETWORK)
	if err != nil {
		return Network{}, err
	}
	defer content.Close()

	var network Network
	if err := json.NewDecoder(content).Decode(&network); err != nil {
		return Network{}, err
	}

	return network, nil
}

func (p PowerDevice) DeviceInformation(ctx context.Context) (Device, error) {
	content, err := p.executeRequest(ctx, DEVICE_INFORMATION)
	if err != nil {
		return Device{}, err
	}
	defer content.Close()

	var device Device
	if err := json.NewDecoder(content).Decode(&device); err != nil {
		return Device{}, err
	}

	return device, nil
}

func (p PowerDevice) Status(ctx context.Context) (PowerStatus, error) {
	content, err := p.executeRequest(ctx, DEVICE_ENERGY)
	if err != nil {
		return PowerStatus{}, err
	}
	defer content.Close()

	var status PowerStatus
	if err := json.NewDecoder(content).Decode(&status); err != nil {
		return PowerStatus{}, err
	}

	return status, nil
}
