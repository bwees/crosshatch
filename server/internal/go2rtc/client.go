// Package go2rtc is an HTTP client for managing camera streams on a go2rtc
// instance (the live-view backend the printers' cameras are published to).
package go2rtc

import (
	"crosshatch/internal/dtos"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"go.uber.org/fx"
)

type Client struct {
	baseURL string
	client  *http.Client
}

func (r *Client) GetStreams() (map[string]dtos.Go2RTCStream, error) {
	res, err := r.client.Get(r.baseURL + "/api/streams")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var streams map[string]dtos.Go2RTCStream
	if err := json.NewDecoder(res.Body).Decode(&streams); err != nil {
		return nil, err
	}

	return streams, nil
}

func (r *Client) DeleteStream(id string) error {
	req, err := http.NewRequest(http.MethodDelete, r.baseURL+"/api/streams/"+id, nil)
	if err != nil {
		return err
	}

	res, err := r.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("failed to delete stream for printer %s: %s", id, res.Status)
	}

	return nil
}

func (r *Client) AddStream(id string, streamURL string) error {
	encodedURL := url.QueryEscape(streamURL)

	req, err := http.NewRequest(
		http.MethodPut,
		fmt.Sprintf("%s/api/streams?name=%s&src=%s", r.baseURL, id, encodedURL),
		nil,
	)
	if err != nil {
		return err
	}

	res, err := r.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("failed to create stream for printer %s: %s", id, res.Status)
	}

	return nil
}

func (r *Client) UpdateStream(id string, streamURL string) error {
	encodedURL := url.QueryEscape(streamURL)
	req, err := http.NewRequest(
		http.MethodPatch,
		fmt.Sprintf("%s/api/streams?name=%s&src=%s", r.baseURL, id, encodedURL),
		nil,
	)
	if err != nil {
		return err
	}

	res, err := r.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("failed to update stream for printer %s: %s", id, res.Status)
	}

	return nil
}

func NewClient() *Client {
	baseURL := os.Getenv("GO2RTC_API_URL")
	if baseURL == "" {
		baseURL = "http://localhost:1984"
	}

	return &Client{
		baseURL: baseURL,
		client:  &http.Client{Timeout: 10 * time.Second},
	}
}

var Module = fx.Provide(NewClient)
