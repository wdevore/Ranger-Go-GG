package engine

import "math"

// LinearEasing is a basic linear lerp
// t,b,c,d
// timePosition ranges from 0 to duration
func LinearEasing(timePosition, startValue, deltaValue, duration float64) float64 {
	return deltaValue*timePosition/duration + startValue
}

// Lerp is basic linear mapping
func Lerp(x1, y1, y float64) float64 {
	x := x1 / y1 * y
	return math.Round(x)
}
