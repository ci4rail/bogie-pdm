package position

import (
	"time"

	"github.com/edgefarm/bogie-pdm/cmd/bogie-edge/internal/gnss"
	"github.com/edgefarm/bogie-pdm/cmd/bogie-edge/internal/steadydrive"
)

// OutputData is the output data of the position unit.
type OutputData struct {
	Timestamp time.Time
	Valid     bool
	Lat       float64
	Lon       float64
	Alt       float64 // m
	Speed     float64 // m/s
}

type inputData struct {
	gnssdata    *gnss.OutputData
	steadydrive *steadydrive.OutputData
}

// Run starts the position unit.
// should be called from a goroutine.
func (p *Unit) Run() {
	p.logger.Info().Msg("running")

	// TODO use steadydrive for sensor fusion
	// TODO use vehicle model to avoid position jumps
	inputCh := p.ps.Sub("gnssraw")

	var inputs inputData
	for {
		select {
		case <-time.After(time.Second):
			p.logger.Debug().Msg("timer")
			var o OutputData
			if inputs.gnssdata != nil {
				// TODO validate gnss data age
				o = OutputData{
					Timestamp: time.Now(),
					Valid:     true,
					Lat:       inputs.gnssdata.Lat,
					Lon:       inputs.gnssdata.Lon,
					Alt:       inputs.gnssdata.Alt,
					Speed:     inputs.gnssdata.Speed,
				}
			} else {
				o = OutputData{
					Timestamp: time.Now(),
					Valid:     false,
				}
			}
			p.ps.Pub(o, "position")

		case msg := <-inputCh:
			switch m := msg.(type) {
			case steadydrive.OutputData:
				inputs.steadydrive = &m
			case gnss.OutputData:
				inputs.gnssdata = &m
			}
		}
	}
}
