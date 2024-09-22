package httpclient_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	httpClient "github.com/vitwit/avail-da-module/relayer/http"
)

func TestHTTPClientHandler_Get(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, world")
	}))
	defer mockServer.Close()

	handler := httpClient.NewHandler()

	resp, err := handler.Get(mockServer.URL)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := "Hello, world\n"
	if string(resp) != expected {
		t.Errorf("expected %s, got %s", expected, string(resp))
	}
}

func TestHTTPClientHandler_GetError(t *testing.T) {
	handler := httpClient.NewHandler()

	_, err := handler.Get("http://invalid-url")
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
}

func TestHTTPClientHandler_Post(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("failed to read request body: %v", err)
		}
		defer r.Body.Close()

		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type application/json, got %s", r.Header.Get("Content-Type"))
		}

		w.Write(body)
	}))
	defer mockServer.Close()

	handler := httpClient.NewHandler()

	jsonData := []byte(`{"name":"test"}`)

	resp, err := handler.Post(mockServer.URL, jsonData)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected := `{"name":"test"}`
	if string(resp) != expected {
		t.Errorf("expected %s, got %s", expected, string(resp))
	}
}

func TestHTTPClientHandler_PostError(t *testing.T) {
	handler := httpClient.NewHandler()

	_, err := handler.Post("http://invalid-url", []byte(`{"data":"test"}`))
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
}
