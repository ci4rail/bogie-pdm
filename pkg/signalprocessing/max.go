package signalprocessing

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
