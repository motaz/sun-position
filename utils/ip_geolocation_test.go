package utils

import (
	"testing"
)

func TestCountryCapitalMap(t *testing.T) {
	// Test that some known countries have their capitals in the map
	testCases := []struct {
		country string
		capital string
	}{
		{"United States", "Washington, D.C."},
		{"United Kingdom", "London"},
		{"France", "Paris"},
		{"Germany", "Berlin"},
		{"Japan", "Tokyo"},
		{"Australia", "Canberra"},
		{"Canada", "Ottawa"},
		{"Sudan", "Khartoum"},
	}

	for _, tc := range testCases {
		t.Run(tc.country, func(t *testing.T) {
			capital, exists := CountryCapitalMap[tc.country]
			if !exists {
				t.Errorf("Country %s not found in CountryCapitalMap", tc.country)
			}
			if capital != tc.capital {
				t.Errorf("Expected capital %s for country %s, but got %s", tc.capital, tc.country, capital)
			}
		})
	}
}

func TestGetCapitalCityForCountry(t *testing.T) {
	// Test that the function returns correct capitals
	testCases := []struct {
		country string
		capital string
	}{
		{"United States", "Washington, D.C."},
		{"United Kingdom", "London"},
		{"France", "Paris"},
		{"Germany", "Berlin"},
		{"Japan", "Tokyo"},
		{"Australia", "Canberra"},
		{"Canada", "Ottawa"},
		{"Sudan", "Khartoum"},
	}

	for _, tc := range testCases {
		t.Run(tc.country, func(t *testing.T) {
			result := GetCapitalCityForCountry(tc.country)
			if result != tc.capital {
				t.Errorf("Expected capital %s for country %s, but got %s", tc.capital, tc.country, result)
			}
		})
	}
}

func TestGetCapitalCityForCountryNotFound(t *testing.T) {
	// Test that the function returns empty string for unknown countries
	result := GetCapitalCityForCountry("NonExistentCountry")
	if result != "" {
		t.Errorf("Expected empty string for non-existent country, but got %s", result)
	}
}

func TestNormalizeCountryName(t *testing.T) {
	// Test that the normalizeCountryName function works correctly
	testCases := []struct {
		input    string
		expected string
	}{
		{"united states", "United States"},
		{"UNITED KINGDOM", "United Kingdom"},
		{"france", "France"},
		{"south korea", "South Korea"},
		{"cape verde", "Cape Verde"},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := normalizeCountryName(tc.input)
			if result != tc.expected {
				t.Errorf("Expected normalized name %s for input %s, but got %s", tc.expected, tc.input, result)
			}
		})
	}
}