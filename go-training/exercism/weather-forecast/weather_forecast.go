// Package weather is script to provide info on current weather conditions.
package weather

// CurrentCondition is an input variable for current weather condition.
var CurrentCondition string

// CurrentLocation is an input variable for current location.
var CurrentLocation string

// Forecast is a function that provides information on current weather condition on a particular location.
func Forecast(city, condition string) string {
	CurrentLocation, CurrentCondition = city, condition
	return CurrentLocation + " - current weather condition: " + CurrentCondition
}
