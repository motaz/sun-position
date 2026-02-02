package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"sun-position/handlers"
)

func TestSunPositionHandlerWithCityName(t *testing.T) {
	// Create a test request with a city name
	req, err := http.NewRequest("GET", "/api/sun-position?city=Khartoum&date=2026-01-28&time=12:00", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler
	handlers.SunPositionHandler(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check if the response body contains expected data
	responseBody := rr.Body.String()
	if !strings.Contains(responseBody, "sun_altitude") {
		t.Errorf("response body does not contain sun_altitude: %s", responseBody)
	}

	if !strings.Contains(responseBody, "sun_azimuth") {
		t.Errorf("response body does not contain sun_azimuth: %s", responseBody)
	}

	// Try to parse the response to validate structure
	var response map[string]interface{}
	err = json.Unmarshal([]byte(responseBody), &response)
	if err != nil {
		t.Fatalf("could not unmarshal response: %v", err)
	}

	// Verify that location contains Khartoum's coordinates
	location, ok := response["location"].(string)
	if !ok {
		t.Errorf("location field not found or not a string in response: %v", response)
	} else if !strings.Contains(location, "15.5007") || !strings.Contains(location, "32.5599") {
		t.Errorf("location does not contain Khartoum's coordinates, got: %s", location)
	}
}

func TestSunPositionHandlerWithNonExistentCity(t *testing.T) {
	// Create a test request with a non-existent city name
	req, err := http.NewRequest("GET", "/api/sun-position?city=NonExistentCity&date=2026-01-28&time=12:00", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler
	handlers.SunPositionHandler(rr, req)

	// Check the status code - should be 400 for bad request
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	// Check if the response body contains the error message
	responseBody := rr.Body.String()
	if !strings.Contains(responseBody, "City not found") {
		t.Errorf("response body does not contain expected error message: %s", responseBody)
	}
}

func TestSunPositionHandlerWithCoordinates(t *testing.T) {
	// Create a test request with coordinates
	req, err := http.NewRequest("GET", "/api/sun-position?lat=40.7128&lon=-74.0060&date=2026-01-28&time=12:00", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler
	handlers.SunPositionHandler(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check if the response body contains expected data
	responseBody := rr.Body.String()
	if !strings.Contains(responseBody, "sun_altitude") {
		t.Errorf("response body does not contain sun_altitude: %s", responseBody)
	}

	if !strings.Contains(responseBody, "sun_azimuth") {
		t.Errorf("response body does not contain sun_azimuth: %s", responseBody)
	}

	// Try to parse the response to validate structure
	var response map[string]interface{}
	err = json.Unmarshal([]byte(responseBody), &response)
	if err != nil {
		t.Fatalf("could not unmarshal response: %v", err)
	}

	// Verify that location contains the provided coordinates
	location, ok := response["location"].(string)
	if !ok {
		t.Errorf("location field not found or not a string in response: %v", response)
	} else if !strings.Contains(location, "40.7128") || !strings.Contains(location, "-74.0060") {
		t.Errorf("location does not contain provided coordinates, got: %s", location)
	}
}

// Mock IP geolocation for testing purposes
func TestSunPositionHandlerWithEmptyCityName(t *testing.T) {
	// Create a test request with no city name, lat, or lon - should default to IP geolocation
	req, err := http.NewRequest("GET", "/api/sun-position?date=2026-01-28&time=12:00", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Add a fake IP to the request
	req.Header.Set("X-Real-IP", "8.8.8.8") // Google's DNS server IP

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler
	handlers.SunPositionHandler(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Logf("Response body: %s", rr.Body.String())
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check if the response body contains expected data
	responseBody := rr.Body.String()
	if !strings.Contains(responseBody, "sun_altitude") {
		t.Errorf("response body does not contain sun_altitude: %s", responseBody)
	}

	if !strings.Contains(responseBody, "sun_azimuth") {
		t.Errorf("response body does not contain sun_azimuth: %s", responseBody)
	}

	// Try to parse the response to validate structure
	var response map[string]interface{}
	err = json.Unmarshal([]byte(responseBody), &response)
	if err != nil {
		t.Fatalf("could not unmarshal response: %v", err)
	}

	// Since we can't predict the exact location from the mock IP,
	// we just verify that the response has the expected fields
	if _, ok := response["sun_altitude"]; !ok {
		t.Errorf("sun_altitude field not found in response: %v", response)
	}

	if _, ok := response["sun_azimuth"]; !ok {
		t.Errorf("sun_azimuth field not found in response: %v", response)
	}

	if _, ok := response["location"]; !ok {
		t.Errorf("location field not found in response: %v", response)
	}
}

