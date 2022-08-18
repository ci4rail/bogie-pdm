package signalprocessing

import (
	"math"
)

// RMS returns the root mean square of the data slice
func RMS(data []float64) float64 {
	var sum float64
	for _, v := range data {
		sum += v * v
	}
	return math.Sqrt(sum / float64(len(data)))
}
