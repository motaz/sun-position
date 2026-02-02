package utils

// City represents a city with its coordinates
type City struct {
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Country   string  `json:"country"`
}

// CommonCities contains a list of common cities with their coordinates
var CommonCities = []City{
	{Name: "Khartoum", Latitude: 15.5007, Longitude: 32.5599, Country: "Sudan"},
	{Name: "Riyadh", Latitude: 24.7136, Longitude: 46.6753, Country: "Saudi Arabia"},
	{Name: "New York", Latitude: 40.7128, Longitude: -74.0060, Country: "USA"},
	{Name: "Los Angeles", Latitude: 34.0522, Longitude: -118.2437, Country: "USA"},
	{Name: "Chicago", Latitude: 41.8781, Longitude: -87.6298, Country: "USA"},
	{Name: "Miami", Latitude: 25.7617, Longitude: -80.1918, Country: "USA"},
	{Name: "London", Latitude: 51.5074, Longitude: -0.1278, Country: "UK"},
	{Name: "Tokyo", Latitude: 35.6762, Longitude: 139.6503, Country: "Japan"},
	{Name: "Paris", Latitude: 48.8566, Longitude: 2.3522, Country: "France"},
	{Name: "Sydney", Latitude: -33.8688, Longitude: 151.2093, Country: "Australia"},
	{Name: "Dubai", Latitude: 25.2048, Longitude: 55.2708, Country: "UAE"},
	{Name: "Singapore", Latitude: 1.3521, Longitude: 103.8198, Country: "Singapore"},
	{Name: "Toronto", Latitude: 43.6532, Longitude: -79.3832, Country: "Canada"},
	{Name: "Berlin", Latitude: 52.5200, Longitude: 13.4050, Country: "Germany"},
	{Name: "Rome", Latitude: 41.9028, Longitude: 12.4964, Country: "Italy"},
	{Name: "Madrid", Latitude: 40.4168, Longitude: -3.7038, Country: "Spain"},
	{Name: "Moscow", Latitude: 55.7558, Longitude: 37.6173, Country: "Russia"},
	{Name: "Beijing", Latitude: 39.9042, Longitude: 116.4074, Country: "China"},
	{Name: "Shanghai", Latitude: 31.2304, Longitude: 121.4737, Country: "China"},
	{Name: "Mumbai", Latitude: 19.0760, Longitude: 72.8777, Country: "India"},
	{Name: "São Paulo", Latitude: -23.5505, Longitude: -46.6333, Country: "Brazil"},
	{Name: "Rio de Janeiro", Latitude: -22.9068, Longitude: -43.1729, Country: "Brazil"},
	{Name: "Mexico City", Latitude: 19.4326, Longitude: -99.1332, Country: "Mexico"},
	{Name: "Cairo", Latitude: 30.0444, Longitude: 31.2357, Country: "Egypt"},
	{Name: "Lagos", Latitude: 6.5244, Longitude: 3.3792, Country: "Nigeria"},
	{Name: "Johannesburg", Latitude: -26.2041, Longitude: 28.0473, Country: "South Africa"},
	{Name: "Seoul", Latitude: 37.5665, Longitude: 126.9780, Country: "South Korea"},
	{Name: "Bangkok", Latitude: 13.7563, Longitude: 100.5018, Country: "Thailand"},
	{Name: "Kuala Lumpur", Latitude: 3.1390, Longitude: 101.6869, Country: "Malaysia"},
	{Name: "Jakarta", Latitude: -6.2088, Longitude: 106.8456, Country: "Indonesia"},
	{Name: "Buenos Aires", Latitude: -34.6037, Longitude: -58.3816, Country: "Argentina"},
	{Name: "Amsterdam", Latitude: 52.3676, Longitude: 4.9041, Country: "Netherlands"},
	{Name: "Vienna", Latitude: 48.2082, Longitude: 16.3738, Country: "Austria"},
	{Name: "Athens", Latitude: 37.9838, Longitude: 23.7275, Country: "Greece"},
	{Name: "Stockholm", Latitude: 59.3293, Longitude: 18.0686, Country: "Sweden"},
	{Name: "Oslo", Latitude: 59.9139, Longitude: 10.7522, Country: "Norway"},
	{Name: "Helsinki", Latitude: 60.1699, Longitude: 24.9384, Country: "Finland"},
	{Name: "Dublin", Latitude: 53.3498, Longitude: -6.2603, Country: "Ireland"},
	{Name: "Brussels", Latitude: 50.8503, Longitude: 4.3517, Country: "Belgium"},
	{Name: "Zurich", Latitude: 47.3769, Longitude: 8.5417, Country: "Switzerland"},
	{Name: "Prague", Latitude: 50.0755, Longitude: 14.4378, Country: "Czech Republic"},
	{Name: "Warsaw", Latitude: 52.2297, Longitude: 21.0122, Country: "Poland"},
	{Name: "Budapest", Latitude: 47.4979, Longitude: 19.0402, Country: "Hungary"},
	{Name: "Lisbon", Latitude: 38.7223, Longitude: -9.1393, Country: "Portugal"},
	{Name: "Copenhagen", Latitude: 55.6761, Longitude: 12.5683, Country: "Denmark"},
	{Name: "Reykjavik", Latitude: 64.1466, Longitude: -21.9426, Country: "Iceland"},
	{Name: "Havana", Latitude: 23.1136, Longitude: -82.3666, Country: "Cuba"},
	{Name: "Kingston", Latitude: 18.1096, Longitude: -77.2975, Country: "Jamaica"},
	{Name: "Panama City", Latitude: 8.9823, Longitude: -79.5199, Country: "Panama"},
	{Name: "Santiago", Latitude: -33.4489, Longitude: -70.6693, Country: "Chile"},
	{Name: "Lima", Latitude: -12.0464, Longitude: -77.0428, Country: "Peru"},
}

// GetCountries returns a list of unique countries from CommonCities
func GetCountries() []string {
	countriesMap := make(map[string]bool)
	var countries []string

	for _, city := range CommonCities {
		if !countriesMap[city.Country] {
			countriesMap[city.Country] = true
			countries = append(countries, city.Country)
		}
	}

	return countries
}

// GetCitiesByCountry returns a list of cities for a given country
func GetCitiesByCountry(country string) []City {
	var cities []City
	for _, city := range CommonCities {
		if city.Country == country {
			cities = append(cities, city)
		}
	}
	return cities
}