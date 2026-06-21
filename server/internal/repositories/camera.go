package repositories

import (
	"crosshatch/internal/dtos"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

type CameraRepository struct {
	baseURL string
}

func (r *CameraRepository) GetStreams() (map[string]dtos.Go2RTCStream, error) {
	res, err := http.Get(r.baseURL + "/api/streams")
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

func (r *CameraRepository) DeleteStream(id string) error {
	req, err := http.NewRequest(http.MethodDelete, r.baseURL+"/api/streams/"+id, nil)
	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	return nil
}

func (r *CameraRepository) AddStream(id string, streamURL string) error {
	encodedURL := url.QueryEscape(streamURL)

	req, err := http.NewRequest(
		http.MethodPut,
		fmt.Sprintf("%s/api/streams?name=%s&src=%s", r.baseURL, id, encodedURL),
		nil,
	)
	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("Failed to create stream for printer %s: %s", id, res.Status)
	}

	return nil
}

func (r *CameraRepository) UpdateStream(id string, streamURL string) error {
	encodedURL := url.QueryEscape(streamURL)
	req, err := http.NewRequest(
		http.MethodPatch,
		fmt.Sprintf("%s/api/streams?name=%s&src=%s", r.baseURL, id, encodedURL),
		nil,
	)
	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("Failed to update stream for printer %s: %s", id, res.Status)
	}

	return nil
}

func NewCameraRepository() *CameraRepository {
	baseURL := os.Getenv("GO2RTC_API_URL")
	if baseURL == "" {
		baseURL = "http://localhost:1984"
	}

	return &CameraRepository{baseURL: baseURL}
}
