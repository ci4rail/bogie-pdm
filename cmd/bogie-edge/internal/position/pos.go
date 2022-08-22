package position

import "time"

// OutputData is the output data of the position unit.
type OutputData struct {
	Timestamp time.Time
	Lat       float64
	Lon       float64
	Alt       float64 // m
	Speed     float64 // m/s
}
