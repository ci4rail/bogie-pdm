package metrics

import (
	"time"

	"github.com/edgefarm/bogie-pdm/cmd/bogie-edge/internal/gnss"
	"github.com/edgefarm/bogie-pdm/cmd/bogie-edge/internal/position"
	"github.com/edgefarm/bogie-pdm/cmd/bogie-edge/internal/steadydrive"
	pb "github.com/edgefarm/bogie-pdm/proto/go/metrics/v1"
	"google.golang.org/protobuf/proto"
)

type metricsData struct {
	steadydrive *steadydrive.OutputData
	position    *position.OutputData
	modem       *modemData
	gnssraw     *gnss.OutputData
}

// Run starts the metrics unit.
// should be called from a goroutine.
func (m *Unit) Run() {
	m.logger.Info().Msg("running")

	inputCh := m.ps.Sub("steadydrive", "position", "cellular", "gnssraw", "internet", "temperature")
	var metrics metricsData

	// start subunits
	go m.runModem()

	for {
		select {
		case <-time.After(time.Duration(m.cfg.PublishPeriod) * time.Second):
			m.logger.Debug().Msg("publish")
			m.publishData(metrics)

		case msg := <-inputCh:
			m.logger.Debug().Msgf("msg %v", msg)
			switch m := msg.(type) {
			case *steadydrive.OutputData:
				metrics.steadydrive = m
			case *position.OutputData:
				metrics.position = m
			case *modemData:
				metrics.modem = m
			case *gnss.OutputData:
				metrics.gnssraw = m
			}
		}
	}
}

func (m *Unit) publishData(metrics metricsData) []byte {
	m.ps.Pub(metrics, "metrics")

	mpb := &pb.Metrics{}
	if metrics.steadydrive != nil {
		mpb.SteadyDrive = &pb.Metrics_SteadyDrive{
			Max: metrics.steadydrive.Max[:],
			Rms: metrics.steadydrive.RMS[:],
		}
	}
	if metrics.position != nil {
		mpb.Position = &pb.Metrics_Position{
			Valid: metrics.position.Valid,
			Lat:   float32(metrics.position.Lat),
			Lon:   float32(metrics.position.Lon),
			Alt:   float32(metrics.position.Alt),
			Speed: float32(metrics.position.Speed),
		}
	}
	// TODO Temperature, Internet
	if metrics.modem != nil {
		mpb.Cellular = &pb.Metrics_Cellular{
			Operator: metrics.modem.OperatorName,
			Strength: metrics.modem.SignalQuality,
		}
	}
	if metrics.gnssraw != nil {
		mpb.GnssRaw = &pb.Metrics_GnssRaw{
			Lat:     float32(metrics.gnssraw.Lat),
			Lon:     float32(metrics.gnssraw.Lon),
			Alt:     float32(metrics.gnssraw.Alt),
			Speed:   float32(metrics.gnssraw.Speed),
			Eph:     float32(metrics.gnssraw.Eph),
			Mode:    int32(metrics.gnssraw.Mode),
			Numsats: int32(metrics.gnssraw.NumSats),
		}
	}

	out, err := proto.Marshal(mpb)
	if err != nil {
		m.logger.Error().Msgf("can't marshall %v", err)
	}
	return out
}
