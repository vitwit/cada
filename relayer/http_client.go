package relayer

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
// Post method
// func (h *HTTPClientHandler) Post(url string, data interface{}) ([]byte, error) {
// 	jsonData, err := json.Marshal(data)
// 	if err != nil {
// 		return nil, fmt.Errorf("error marshalling data: %v", err)
// 	}

// 	resp, err := h.client.Post(url, "application/json", bytes.NewBuffer(jsonData))
// 	if err != nil {
// 		return nil, fmt.Errorf("POST request error: %v", err)
// 	}
// 	defer resp.Body.Close()

// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, fmt.Errorf("error reading response body: %v", err)
// 	}
// 	return body, nil
// }

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

// RequestData struct for the POST request payload
type RequestData struct {
	Data      []byte `json:"data,omitempty"`
	Extrinsic string `json:"extrinsic,omitempty"`
}

type BlockInfo struct {
	BlockNumber int    `json:"block_number"`
	BlockHash   string `json:"block_hash"`
	Hash        string `json:"hash"`
	Index       int    `json:"index"`
}
