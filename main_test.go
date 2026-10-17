package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestHealthHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(healthHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := "ok"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestStatusAPIHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/status", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(statusAPIHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check Content-Type
	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("handler returned wrong content type: got %v want %v",
			contentType, "application/json")
	}

	// Check JSON structure
	var status StatusInfo
	err = json.NewDecoder(rr.Body).Decode(&status)
	if err != nil {
		t.Errorf("failed to decode JSON response: %v", err)
	}

	// Verify required fields are present
	if status.Status == "" {
		t.Error("status field is empty")
	}
	if status.Version == "" {
		t.Error("version field is empty")
	}
}

func TestMetricsHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/metrics", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(metricsHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check Content-Type
	contentType := rr.Header().Get("Content-Type")
	if contentType != "text/plain" {
		t.Errorf("handler returned wrong content type: got %v want %v",
			contentType, "text/plain")
	}

	body := rr.Body.String()
	if len(body) == 0 {
		t.Error("metrics response is empty")
	}
}

func TestRootHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(rootHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check Content-Type
	contentType := rr.Header().Get("Content-Type")
	if contentType != "text/html" {
		t.Errorf("handler returned wrong content type: got %v want %v",
			contentType, "text/html")
	}

	body := rr.Body.String()
	if len(body) == 0 {
		t.Error("HTML response is empty")
	}
}

func TestGetEnv(t *testing.T) {
	// Test with existing env var
	os.Setenv("TEST_VAR", "test_value")
	defer os.Unsetenv("TEST_VAR")

	result := getEnv("TEST_VAR", "fallback")
	if result != "test_value" {
		t.Errorf("getEnv returned wrong value: got %v want %v",
			result, "test_value")
	}

	// Test with fallback
	result = getEnv("NONEXISTENT_VAR", "fallback")
	if result != "fallback" {
		t.Errorf("getEnv returned wrong value: got %v want %v",
			result, "fallback")
	}
}

func TestGetStatusInfo(t *testing.T) {
	os.Setenv("ENVIRONMENT", "test")
	os.Setenv("POD_NAME", "test-pod")
	os.Setenv("POD_NAMESPACE", "test-namespace")
	defer func() {
		os.Unsetenv("ENVIRONMENT")
		os.Unsetenv("POD_NAME")
		os.Unsetenv("POD_NAMESPACE")
	}()

	status := getStatusInfo()

	if status.Environment != "test" {
		t.Errorf("Expected environment 'test', got '%s'", status.Environment)
	}
	if status.PodName != "test-pod" {
		t.Errorf("Expected pod name 'test-pod', got '%s'", status.PodName)
	}
	if status.Namespace != "test-namespace" {
		t.Errorf("Expected namespace 'test-namespace', got '%s'", status.Namespace)
	}
	if status.Status != "healthy" {
		t.Errorf("Expected status 'healthy', got '%s'", status.Status)
	}
}
