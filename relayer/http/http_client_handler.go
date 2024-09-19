package http_client

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

// HTTPClientHandler struct
type HTTPClientHandler struct {
	client *http.Client
}

// NewHTTPClientHandler creates a new HTTPClientHandler with default settings
func NewHTTPClientHandler() *HTTPClientHandler {
	return &HTTPClientHandler{
		client: &http.Client{
			Timeout: 100 * time.Second,
		},
	}
}

// Get method
func (h *HTTPClientHandler) Get(url string) ([]byte, error) {
	resp, err := h.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("GET request error: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}
	return body, nil
}

// Post method
func (h *HTTPClientHandler) Post(url string, jsonData []byte) ([]byte, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("POST request error: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}
	return body, nil
}
