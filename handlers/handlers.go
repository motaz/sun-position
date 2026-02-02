package handlers

import (
	"encoding/json"
	"fmt"
	"math"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"sun-position/utils"
)

// HomeHandler serves the main page
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(indexHTML)
}

// SunPositionRequest represents the request body for sun position calculation
type SunPositionRequest struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Date      string  `json:"date"` // Format: YYYY-MM-DD
	Time      string  `json:"time"` // Format: HH:MM
}

// SunPositionResponse represents the response body
type SunPositionResponse struct {
	SunAltitude float64   `json:"sun_altitude"`
	SunAzimuth  float64   `json:"sun_azimuth"`
	Timestamp   time.Time `json:"timestamp"`
	Location    string    `json:"location"`
	City        string    `json:"city,omitempty"`  // Include city name if available
	Date        string    `json:"date"`
	Time        string    `json:"time"`
	Sunrise     string    `json:"sunrise"`
	Sunset      string    `json:"sunset"`
}


// getClientIP extracts the client's IP address from the request, considering proxies
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header (first entry if multiple IPs)
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		ip := strings.Split(forwarded, ",")[0]
		return strings.TrimSpace(ip)
	}

	// Check X-Real-IP header
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// Use RemoteAddr as fallback
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

// SunPositionHandler calculates and returns the sun's position
func SunPositionHandler(w http.ResponseWriter, r *http.Request) {
	var req SunPositionRequest

	// Parse query parameters
	cityName := r.URL.Query().Get("city")
	latStr := r.URL.Query().Get("lat")
	lonStr := r.URL.Query().Get("lon")

	var lat, lon float64
	var err error

	// If city name is provided, use it to get coordinates
	if cityName != "" {
		cityFound := false
		for _, city := range utils.CommonCities {
			// Check if the city name matches exactly or partially
			if strings.EqualFold(city.Name, cityName) {
				lat = city.Latitude
				lon = city.Longitude
				cityFound = true
				break
			}
		}
		if !cityFound {
			http.Error(w, "City not found", http.StatusBadRequest)
			return
		}
	} else {
		// If no city is provided, try to get location from IP address
		if latStr == "" && lonStr == "" {
			// Get client IP address
			clientIP := getClientIP(r)

			// Attempt to get location from IP
			location, err := utils.GetLocationFromIP(clientIP)
			if err == nil && location != nil && location.Country != "" {
				// Get the capital city for the detected country
				capitalCity := utils.GetCapitalCityForCountry(location.Country)

				// Look for the capital city in our CommonCities list
				if capitalCity != "" {
					for _, city := range utils.CommonCities {
						if strings.EqualFold(city.Name, capitalCity) {
							lat = city.Latitude
							lon = city.Longitude
							cityName = capitalCity // Update cityName to the detected capital
							break
						}
					}
				}

				// If capital city is not in our list or wasn't found, use coordinates from IP if available
				if lat == 0 && lon == 0 && location.Lat != "" && location.Lon != "" {
					parsedLat, err1 := strconv.ParseFloat(location.Lat, 64)
					parsedLon, err2 := strconv.ParseFloat(location.Lon, 64)
					if err1 == nil && err2 == nil {
						lat = parsedLat
						lon = parsedLon
						cityName = capitalCity // Use capital city name if available, otherwise it remains empty
					}
				}
			}

			// If we still don't have coordinates, default to Khartoum
			if lat == 0 && lon == 0 {
				lat = 15.5007
				lon = 32.5599
				if cityName == "" {
					cityName = "Khartoum"
				}
			}
		} else {
			// Parse latitude and longitude from query params
			lat, err = strconv.ParseFloat(latStr, 64)
			if err != nil {
				http.Error(w, "Invalid latitude", http.StatusBadRequest)
				return
			}

			lon, err = strconv.ParseFloat(lonStr, 64)
			if err != nil {
				http.Error(w, "Invalid longitude", http.StatusBadRequest)
				return
			}
		}
	}

	// Use current date/time if not provided
	dateStr := r.URL.Query().Get("date")
	timeStr := r.URL.Query().Get("time")

	if dateStr == "" || timeStr == "" {
		now := time.Now()
		dateStr = now.Format("2006-01-02")
		timeStr = now.Format("15:04")
	}

	req.Latitude = lat
	req.Longitude = lon
	req.Date = dateStr
	req.Time = timeStr

	// Parse the date and time
	// The user enters local time for the location, so we need to interpret it correctly
	// Calculate the time zone offset based on longitude (each 15 degrees = 1 hour)
	// Round to the nearest time zone (multiples of 15 degrees)
	timeZoneOffsetHours := math.Round(req.Longitude / 15.0)
	timeZoneOffsetSeconds := int(timeZoneOffsetHours * 60 * 60) // Convert hours to seconds
	location := time.FixedZone("Local", timeZoneOffsetSeconds)

	parsedTime, err := time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", req.Date, req.Time), location)
	if err != nil {
		http.Error(w, "Invalid date or time format", http.StatusBadRequest)
		return
	}

	// Debug: Print the parsed time
	fmt.Printf("DEBUG: Input time %s %s for location (%.4f, %.4f) -> parsed as %v (location: %s)\n",
		req.Date, req.Time, req.Latitude, req.Longitude, parsedTime, location.String())

	// Calculate sun position
	altitude, azimuth := utils.CalculateSunPosition(req.Latitude, req.Longitude, parsedTime)

	// Determine the city name to include in the response
	var responseCityName string
	if cityName != "" {
		responseCityName = cityName
	} else if req.Latitude == 15.5007 && req.Longitude == 32.5599 {
		// If using default coordinates, set city name to Khartoum
		responseCityName = "Khartoum"
	} else {
		// Try to find the city name from the CommonCities list
		for _, city := range utils.CommonCities {
			if math.Abs(city.Latitude-req.Latitude) < 0.01 && math.Abs(city.Longitude-req.Longitude) < 0.01 {
				responseCityName = city.Name
				break
			}
		}
	}

	response := SunPositionResponse{
		SunAltitude: altitude,
		SunAzimuth:  azimuth,
		Timestamp:   parsedTime,
		Location:    fmt.Sprintf("%.4f, %.4f", req.Latitude, req.Longitude),
		City:        responseCityName, // Include city name in response
		Date:        req.Date,
		Time:        req.Time,
	}

	// Calculate sunrise and sunset for the given date and location
	sunriseTime, sunsetTime := utils.CalculateSunriseSunset(req.Latitude, req.Longitude, parsedTime)
	if !sunriseTime.IsZero() {
		response.Sunrise = sunriseTime.Format("15:04")
	} else {
		response.Sunrise = "N/A"
	}
	if !sunsetTime.IsZero() {
		response.Sunset = sunsetTime.Format("15:04")
	} else {
		response.Sunset = "N/A"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
