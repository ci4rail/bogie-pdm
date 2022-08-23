package metrics

import (
	"time"

	"github.com/edgefarm/bogie-pdm/cmd/bogie-edge/internal/position"
	"github.com/edgefarm/bogie-pdm/cmd/bogie-edge/internal/steadydrive"
)

type metricsData struct {
	steadydrive *steadydrive.OutputData
	position    *position.OutputData
}

// Run starts the metrics unit.
// should be called from a goroutine.
func (m *Unit) Run() {
	m.logger.Info().Msg("running")

	inputCh := m.ps.Sub("steadydrive", "position", "cellular", "gnss", "internet", "temperature")
	var metrics metricsData

	for {
		select {
		case <-time.After(time.Duration(m.cfg.PublishPeriod) * time.Second):
			m.logger.Debug().Msg("publish")
		case msg := <-inputCh:
			m.logger.Debug().Msgf("msg %v", msg)
			switch m := msg.(type) {
			case *steadydrive.OutputData:
				metrics.steadydrive = m
			case *position.OutputData:
				metrics.position = m
			}
		}
	}
}
