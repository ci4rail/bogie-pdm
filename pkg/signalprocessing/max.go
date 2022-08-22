package signalprocessing

import "math"

// Max returns the maximum value of the data slice
func Max(data []float64) float64 {
	var max float64
	for _, v := range data {
		if v > max {
			max = v
		}
	}
	return max
}

// MaxAbs returns the maximum value of the absolute values of the data slice
func MaxAbs(data []float64) float64 {
	var max float64
	for _, v := range data {
		a := math.Abs(v)

		if a > max {
			max = a
		}
	}
	return max
}
