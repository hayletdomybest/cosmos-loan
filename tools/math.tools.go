package tools

import "math"

// Round rounds a float64 to the specified number of decimal places.
func Round(value float64, places int) float64 {
	factor := math.Pow(10, float64(places))
	return math.Round(value*factor) / factor
}

// Ceil rounds up a float64 to the specified number of decimal places.
func Ceil(value float64, places int) float64 {
	factor := math.Pow(10, float64(places))
	return math.Ceil(value*factor) / factor
}
