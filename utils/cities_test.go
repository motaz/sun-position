package utils

import (
	"testing"
)

func TestCommonCitiesContainsKhartoum(t *testing.T) {
	found := false
	expectedCity := "Khartoum"
	expectedCountry := "Sudan"
	expectedLat := 15.5007
	expectedLng := 32.5599

	for _, city := range CommonCities {
		if city.Name == expectedCity && city.Country == expectedCountry {
			if city.Latitude == expectedLat && city.Longitude == expectedLng {
				found = true
				break
			}
		}
	}

	if !found {
		t.Errorf("Expected to find Khartoum in CommonCities, but it was not found")
	}
}

func TestCommonCitiesNotEmpty(t *testing.T) {
	if len(CommonCities) == 0 {
		t.Errorf("Expected CommonCities to not be empty, but it has %d cities", len(CommonCities))
	}

	if len(CommonCities) < 10 { // We know we have many more cities
		t.Errorf("Expected CommonCities to have more cities, but it only has %d", len(CommonCities))
	}
}

func TestCommonCitiesHasValidCoordinates(t *testing.T) {
	for i, city := range CommonCities {
		if city.Name == "" {
			t.Errorf("City at index %d has an empty name", i)
		}

		if city.Country == "" {
			t.Errorf("City %s has an empty country", city.Name)
		}

		// Latitude should be between -90 and 90
		if city.Latitude < -90 || city.Latitude > 90 {
			t.Errorf("City %s has invalid latitude: %f", city.Name, city.Latitude)
		}

		// Longitude should be between -180 and 180
		if city.Longitude < -180 || city.Longitude > 180 {
			t.Errorf("City %s has invalid longitude: %f", city.Name, city.Longitude)
		}
	}
}

func TestFindCityByName(t *testing.T) {
	testCases := []struct {
		name           string
		expectedFound  bool
	}{
		{"Khartoum", true},
		{"Riyadh", true},
		{"London", true},
		{"Tokyo", true},
		{"New York", true},
		{"NonExistentCity", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			found := false
			for _, city := range CommonCities {
				if city.Name == tc.name {
					found = true
					break
				}
			}

			if found != tc.expectedFound {
				t.Errorf("Expected to find %s: %t, but got: %t", tc.name, tc.expectedFound, found)
			}
		})
	}
}

