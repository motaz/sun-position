package utils

import (
	"math"
	"time"
)

// CalculateSunPosition calculates the sun's altitude and azimuth for a given location and time
func CalculateSunPosition(latitude, longitude float64, dateTime time.Time) (altitude, azimuth float64) {
	// Convert degrees to radians
	latRad := latitude * math.Pi / 180

	// Calculate day of year
	year, month, day := dateTime.Date()
	dayOfYear := daysSinceJan1(year, month, day)

	// Calculate declination angle (δ) - more accurate formula
	declination := calculateDeclinationAccurate(dayOfYear)
	declinationRad := declination

	// Calculate equation of time (in minutes) - more accurate formula
	equationOfTime := calculateEquationOfTimeAccurate(dayOfYear)

	// Calculate solar time
	solarTime := calculateSolarTimeAccurate(dateTime, longitude, equationOfTime)

	// Calculate hour angle (H) in radians
	// Hour angle is 0 at solar noon, negative in the morning, positive in the afternoon
	// Each hour corresponds to 15 degrees (360/24)
	hourAngleDeg := (solarTime - 12.0) * 15.0
	hourAngle := hourAngleDeg * math.Pi / 180.0

	// Calculate solar altitude (α_s)
	sinAltitude := math.Sin(latRad)*math.Sin(declinationRad) +
		math.Cos(latRad)*math.Cos(declinationRad)*math.Cos(hourAngle)
	altitude = math.Asin(sinAltitude) * 180 / math.Pi

	// Apply atmospheric refraction correction for low altitudes
	// This correction is most significant when the sun is near the horizon
	altitudeRad := altitude * math.Pi / 180.0
	refractionCorrection := 0.0
	if altitude > -0.575 {
		// Formula for atmospheric refraction (in degrees)
		// Valid for altitudes above -0.575 degrees
		refractionCorrection = 0.016667 / math.Tan(altitudeRad + 0.003138/(altitudeRad + 0.089186))
	} else {
		// For altitudes below -0.575 degrees, use a different approximation
		refractionCorrection = 0.57644 * math.Exp(-0.00149*altitude) - 0.07156
	}

	// Apply the refraction correction
	altitude += refractionCorrection

	// Calculate solar azimuth (γ_s) - more accurate calculation
	// Using the formula from Reda and Andreas (2005)
	cosAzimuth := (math.Sin(declinationRad)*math.Cos(latRad) -
		math.Cos(declinationRad)*math.Sin(latRad)*math.Cos(hourAngle)) /
		math.Cos(altitude*math.Pi/180)

	// Clamp cosine value to [-1, 1] range to avoid domain errors
	if cosAzimuth > 1 {
		cosAzimuth = 1
	} else if cosAzimuth < -1 {
		cosAzimuth = -1
	}

	sinAzimuth := math.Sin(hourAngle) * math.Cos(declinationRad) / math.Cos(altitude*math.Pi/180)

	azimuth = math.Atan2(sinAzimuth, cosAzimuth) * 180 / math.Pi

	// Ensure azimuth is in 0-360 range
	if azimuth < 0 {
		azimuth += 360
	}

	return altitude, azimuth
}

// calculateDeclinationAccurate calculates the solar declination angle in radians (more accurate)
func calculateDeclinationAccurate(dayOfYear int) float64 {
	// More accurate formula for solar declination
	gamma := 2 * math.Pi * float64(dayOfYear-1) / 365.0
	declination := 0.006918 -
		0.399912*math.Cos(gamma) +
		0.070257*math.Sin(gamma) -
		0.006758*math.Cos(2*gamma) +
		0.000907*math.Sin(2*gamma) -
		0.002697*math.Cos(3*gamma) +
		0.00148*math.Sin(3*gamma)

	return declination
}

// calculateEquationOfTimeAccurate calculates the equation of time in minutes (more accurate)
func calculateEquationOfTimeAccurate(dayOfYear int) float64 {
	gamma := 2 * math.Pi * float64(dayOfYear-1) / 365.0

	// Equation of time in minutes
	eot := 229.18 * (0.000075 +
		0.001868*math.Cos(gamma) -
		0.032077*math.Sin(gamma) -
		0.014615*math.Cos(2*gamma) -
		0.040849*math.Sin(2*gamma))

	return eot
}

// calculateSolarTimeAccurate calculates the solar time in hours (more accurate)
func calculateSolarTimeAccurate(dateTime time.Time, longitude float64, equationOfTime float64) float64 {
	// Standard time in decimal hours
	hour, min, sec := dateTime.Clock()
	standardTime := float64(hour) + float64(min)/60.0 + float64(sec)/3600.0

	// Time correction factor (minutes)
	// 4 minutes per degree difference between longitude and the standard meridian for the time zone
	// Plus the equation of time correction
	// Standard longitude for the time zone (multiples of 15 degrees from Greenwich)
	// Using round instead of floor to get the closest standard meridian
	standardMeridian := math.Round(longitude / 15.0) * 15.0

	// 4 minutes per degree longitude difference from standard meridian
	longitudeCorrection := 4.0 * (longitude - standardMeridian)

	// Total time correction in minutes
	timeCorrection := equationOfTime + longitudeCorrection

	// Apply correction to get solar time
	correctedTime := standardTime + timeCorrection/60.0

	// Normalize solar time to 24-hour format
	for correctedTime < 0 {
		correctedTime += 24
	}
	for correctedTime >= 24 {
		correctedTime -= 24
	}

	return correctedTime
}

// calculateHourAngleAccurate calculates the hour angle in radians (more accurate)
func calculateHourAngleAccurate(solarTime float64) float64 {
	// Hour angle in degrees (15 degrees per hour from solar noon)
	// Solar noon is at 12:00 solar time
	hourAngleDeg := (solarTime - 12) * 15
	// Convert to radians
	return hourAngleDeg * math.Pi / 180
}

// daysSinceJan1 calculates the day of year (1-365/366)
func daysSinceJan1(year int, month time.Month, day int) int {
	// Days in each month (non-leap year)
	daysInMonth := []int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}

	// Check for leap year
	if isLeapYear(year) {
		daysInMonth[1] = 29
	}

	totalDays := 0
	for i := 0; i < int(month)-1; i++ {
		totalDays += daysInMonth[i]
	}
	totalDays += day

	return totalDays
}

// isLeapYear checks if a year is a leap year
func isLeapYear(year int) bool {
	return (year%4 == 0) && (year%100 != 0) || (year%400 == 0)
}

// CalculateSunriseSunset computes approximate sunrise and sunset times for the given date and location.
// Returns zero times when sunrise or sunset cannot be determined (e.g., polar day/night).
func CalculateSunriseSunset(latitude, longitude float64, date time.Time) (time.Time, time.Time) {
	// NOAA-based sunrise/sunset calculation (approximate)
	// Reference: https://gml.noaa.gov/grad/solcalc/solareqns.PDF (simplified)
	year, month, day := date.Date()
	dayOfYear := daysSinceJan1(year, month, day)

	// Solar declination (radians) and equation of time (minutes)
	decl := calculateDeclinationAccurate(dayOfYear)   // radians
	eot := calculateEquationOfTimeAccurate(dayOfYear) // minutes

	// Convert latitude to radians
	latRad := latitude * math.Pi / 180.0

	// Sun altitude for sunrise/sunset including refraction
	h0 := -0.833 * math.Pi / 180.0

	// Calculate the hour angle H0 (radians)
	cosH0 := (math.Sin(h0) - math.Sin(latRad)*math.Sin(decl)) / (math.Cos(latRad) * math.Cos(decl))
	if cosH0 > 1 || cosH0 < -1 {
		// Sun does not rise or does not set on this date at this location
		return time.Time{}, time.Time{}
	}
	H0 := math.Acos(cosH0) // radians

	// Solar noon in UTC (hours)
	// solarNoonUTC = 12 - (longitude / 15) - (EoT / 60)
	solarNoonUTC := 12.0 - (longitude / 15.0) - (eot / 60.0)

	// Convert hour angle to hours: H0 (radians) -> hours = H0 * 12 / pi
	deltaHours := H0 * 12.0 / math.Pi

	sunriseUTC := solarNoonUTC - deltaHours
	sunsetUTC := solarNoonUTC + deltaHours

	// Convert UTC hour values to time.Time at UTC
	locUTC := time.FixedZone("UTC", 0)
	startOfDayUTC := time.Date(year, month, day, 0, 0, 0, 0, locUTC)

	sunriseTimeUTC := startOfDayUTC.Add(time.Duration(sunriseUTC * float64(time.Hour)))
	sunsetTimeUTC := startOfDayUTC.Add(time.Duration(sunsetUTC * float64(time.Hour)))

	// Approximate local timezone offset from longitude (hours)
	tzOffsetHours := int(math.Round(longitude / 15.0))
	localLoc := time.FixedZone("Local", tzOffsetHours*3600)

	// Return times converted to approximate local timezone
	return sunriseTimeUTC.In(localLoc), sunsetTimeUTC.In(localLoc)
}
