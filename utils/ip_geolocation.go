package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// CountryCapitalMap maps country names to their capital cities
var CountryCapitalMap = map[string]string{
	"Afghanistan":         "Kabul",
	"Albania":             "Tirana",
	"Algeria":             "Algiers",
	"American Samoa":      "Pago Pago",
	"Andorra":             "Andorra la Vella",
	"Angola":              "Luanda",
	"Antigua and Barbuda": "St. John's",
	"Argentina":           "Buenos Aires",
	"Armenia":             "Yerevan",
	"Australia":           "Canberra",
	"Austria":             "Vienna",
	"Azerbaijan":          "Baku",
	"Bahamas":             "Nassau",
	"Bahrain":             "Manama",
	"Bangladesh":          "Dhaka",
	"Barbados":            "Bridgetown",
	"Belarus":             "Minsk",
	"Belgium":             "Brussels",
	"Belize":              "Belmopan",
	"Benin":               "Porto-Novo",
	"Bhutan":              "Thimphu",
	"Bolivia":             "La Paz",
	"Bosnia and Herzegovina": "Sarajevo",
	"Botswana":            "Gaborone",
	"Brazil":              "Brasília",
	"Brunei":              "Bandar Seri Begawan",
	"Bulgaria":            "Sofia",
	"Burkina Faso":        "Ouagadougou",
	"Burundi":             "Gitega",
	"Cabo Verde":          "Praia",
	"Cambodia":            "Phnom Penh",
	"Cameroon":            "Yaoundé",
	"Canada":              "Ottawa",
	"Central African Republic": "Bangui",
	"Chad":                "N'Djamena",
	"Chile":               "Santiago",
	"China":               "Beijing",
	"Colombia":            "Bogotá",
	"Comoros":             "Moroni",
	"Congo, Democratic Republic of the": "Kinshasa",
	"Congo, Republic of the": "Brazzaville",
	"Costa Rica":          "San José",
	"Croatia":             "Zagreb",
	"Cuba":                "Havana",
	"Cyprus":              "Nicosia",
	"Czech Republic":      "Prague",
	"Denmark":             "Copenhagen",
	"Djibouti":            "Djibouti",
	"Dominica":            "Roseau",
	"Dominican Republic":  "Santo Domingo",
	"Ecuador":             "Quito",
	"Egypt":               "Cairo",
	"El Salvador":         "San Salvador",
	"Equatorial Guinea":   "Malabo",
	"Eritrea":             "Asmara",
	"Estonia":             "Tallinn",
	"Eswatini":            "Mbabane",
	"Ethiopia":            "Addis Ababa",
	"Fiji":                "Suva",
	"Finland":             "Helsinki",
	"France":              "Paris",
	"Gabon":               "Libreville",
	"Gambia":              "Banjul",
	"Georgia":             "Tbilisi",
	"Germany":             "Berlin",
	"Ghana":               "Accra",
	"Greece":              "Athens",
	"Grenada":             "St. George's",
	"Guatemala":           "Guatemala City",
	"Guinea":              "Conakry",
	"Guinea-Bissau":       "Bissau",
	"Guyana":              "Georgetown",
	"Haiti":               "Port-au-Prince",
	"Honduras":            "Tegucigalpa",
	"Hungary":             "Budapest",
	"Iceland":             "Reykjavik",
	"India":               "New Delhi",
	"Indonesia":           "Jakarta",
	"Iran":                "Tehran",
	"Iraq":                "Baghdad",
	"Ireland":             "Dublin",
	"Israel":              "Jerusalem",
	"Italy":               "Rome",
	"Ivory Coast":         "Yamoussoukro",
	"Jamaica":             "Kingston",
	"Japan":               "Tokyo",
	"Jordan":              "Amman",
	"Kazakhstan":          "Astana",
	"Kenya":               "Nairobi",
	"Kiribati":            "Tarawa",
	"Korea, North":        "Pyongyang",
	"Korea, South":        "Seoul",
	"Kosovo":              "Pristina",
	"Kuwait":              "Kuwait City",
	"Kyrgyzstan":          "Bishkek",
	"Laos":                "Vientiane",
	"Latvia":              "Riga",
	"Lebanon":             "Beirut",
	"Lesotho":             "Maseru",
	"Liberia":             "Monrovia",
	"Libya":               "Tripoli",
	"Liechtenstein":       "Vaduz",
	"Lithuania":           "Vilnius",
	"Luxembourg":          "Luxembourg City",
	"Madagascar":          "Antananarivo",
	"Malawi":              "Lilongwe",
	"Malaysia":            "Kuala Lumpur",
	"Maldives":            "Malé",
	"Mali":                "Bamako",
	"Malta":               "Valletta",
	"Marshall Islands":    "Majuro",
	"Mauritania":          "Nouakchott",
	"Mauritius":           "Port Louis",
	"Mexico":              "Mexico City",
	"Micronesia":          "Palikir",
	"Moldova":             "Chișinău",
	"Monaco":              "Monaco",
	"Mongolia":            "Ulaanbaatar",
	"Montenegro":          "Podgorica",
	"Morocco":             "Rabat",
	"Mozambique":          "Maputo",
	"Myanmar":             "Naypyidaw",
	"Namibia":             "Windhoek",
	"Nauru":               "Yaren",
	"Nepal":               "Kathmandu",
	"Netherlands":         "Amsterdam",
	"New Zealand":         "Wellington",
	"Nicaragua":           "Managua",
	"Niger":               "Niamey",
	"Nigeria":             "Abuja",
	"North Macedonia":     "Skopje",
	"Norway":              "Oslo",
	"Oman":                "Muscat",
	"Pakistan":            "Islamabad",
	"Palau":               "Ngerulmud",
	"Panama":              "Panama City",
	"Papua New Guinea":    "Port Moresby",
	"Paraguay":            "Asunción",
	"Peru":                "Lima",
	"Philippines":         "Manila",
	"Poland":              "Warsaw",
	"Portugal":            "Lisbon",
	"Qatar":               "Doha",
	"Romania":             "Bucharest",
	"Russia":              "Moscow",
	"Rwanda":              "Kigali",
	"Saint Kitts and Nevis": "Basseterre",
	"Saint Lucia":         "Castries",
	"Saint Vincent and the Grenadines": "Kingstown",
	"Samoa":               "Apia",
	"San Marino":          "San Marino",
	"Sao Tome and Principe": "São Tomé",
	"Saudi Arabia":        "Riyadh",
	"Senegal":             "Dakar",
	"Serbia":              "Belgrade",
	"Seychelles":          "Victoria",
	"Sierra Leone":        "Freetown",
	"Singapore":           "Singapore",
	"Slovakia":            "Bratislava",
	"Slovenia":            "Ljubljana",
	"Solomon Islands":     "Honiara",
	"Somalia":             "Mogadishu",
	"South Africa":        "Pretoria",
	"South Sudan":         "Juba",
	"Spain":               "Madrid",
	"Sri Lanka":           "Colombo",
	"Sudan":               "Khartoum",
	"Suriname":            "Paramaribo",
	"Sweden":              "Stockholm",
	"Switzerland":         "Bern",
	"Syria":               "Damascus",
	"Taiwan":              "Taipei",
	"Tajikistan":          "Dushanbe",
	"Tanzania":            "Dodoma",
	"Thailand":            "Bangkok",
	"Timor-Leste":         "Dili",
	"Togo":                "Lomé",
	"Tonga":               "Nuku'alofa",
	"Trinidad and Tobago": "Port of Spain",
	"Tunisia":             "Tunis",
	"Turkey":              "Ankara",
	"Turkmenistan":        "Ashgabat",
	"Tuvalu":              "Funafuti",
	"Uganda":              "Kampala",
	"Ukraine":             "Kyiv",
	"United Arab Emirates": "Abu Dhabi",
	"United Kingdom":      "London",
	"United States":       "Washington, D.C.",
	"Uruguay":             "Montevideo",
	"Uzbekistan":          "Tashkent",
	"Vanuatu":             "Port Vila",
	"Vatican City":        "Vatican City",
	"Venezuela":           "Caracas",
	"Vietnam":             "Hanoi",
	"Yemen":               "Sana'a",
	"Zambia":              "Lusaka",
	"Zimbabwe":            "Harare",
}

// IPLocation represents the response from the IP geolocation service
type IPLocation struct {
	Country string `json:"country"`
	City    string `json:"city"`
	Region  string `json:"region"`
	Lat     string `json:"lat"`
	Lon     string `json:"lon"`
}

// GetLocationFromIP gets the location information from an IP address using ipapi.co
func GetLocationFromIP(ip string) (*IPLocation, error) {
	// If IP is empty or localhost, return a default location
	if ip == "" || ip == "127.0.0.1" || ip == "::1" {
		ip = ""
	}

	url := fmt.Sprintf("https://ipapi.co/%s/json/", ip)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var location IPLocation
	err = json.Unmarshal(body, &location)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &location, nil
}

// GetCapitalCityForCountry returns the capital city for a given country
func GetCapitalCityForCountry(country string) string {
	// Normalize the country name to title case for comparison
	countryNormalized := normalizeCountryName(country)
	
	if capital, exists := CountryCapitalMap[countryNormalized]; exists {
		return capital
	}
	
	// If not found, try to find by partial match
	for countryName, capital := range CountryCapitalMap {
		if strings.Contains(strings.ToLower(countryName), strings.ToLower(country)) ||
		   strings.Contains(strings.ToLower(country), strings.ToLower(countryName)) {
			return capital
		}
	}
	
	return ""
}

// normalizeCountryName normalizes the country name to title case
func normalizeCountryName(name string) string {
	words := strings.Split(name, " ")
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.Title(strings.ToLower(word))
		}
	}
	return strings.Join(words, " ")
}