package metrics

import (
	"sync"
	"time"

	"github.com/ci4rail/bogie-pdm/cmd/bogie-edge/internal/gnss"
	"github.com/ci4rail/bogie-pdm/cmd/bogie-edge/internal/position"
	"github.com/ci4rail/bogie-pdm/cmd/bogie-edge/internal/steadydrive"
	pb "github.com/ci4rail/bogie-pdm/proto/go/metrics/v1"
	"google.golang.org/protobuf/proto"
)

type metricsData struct {
	mutex       *sync.Mutex
	steadydrive *steadydrive.OutputData
	position    *position.OutputData
	modem       *modemData
	gnssraw     *gnss.OutputData
	temp        *temperatureData
	internet    *internetData
}

// Run starts the metrics unit.
// should be called from a goroutine.
func (m *Unit) Run() {
	m.logger.Info().Msg("running")

	inputCh := m.ps.Sub("steadydrive", "position", "cellular", "gnssraw", "internet", "temperature")
	var metr *metricsData = &metricsData{
		mutex: &sync.Mutex{},
	}

	// start subunits
	go m.runModem()
	go m.runTemperature()
	go m.runInternet()

	// publish go routine
	go func(metr *metricsData) {
		for {
			time.Sleep(time.Duration(m.cfg.PublishPeriod) * time.Second)
			metr.mutex.Lock()
			o := m.publishData(metr)
			err := m.export.PubExport("metrics", o)
			if err != nil {
				m.logger.Error().Msgf("can't publish %v", err)
			}
			metr.mutex.Unlock()
		}
	}(metr)

	for msg := range inputCh {
		metr.mutex.Lock()
		switch m := msg.(type) {
		case steadydrive.OutputData:
			metr.steadydrive = &m
		case position.OutputData:
			metr.position = &m
		case modemData:
			metr.modem = &m
		case temperatureData:
			metr.temp = &m
		case internetData:
			metr.internet = &m
		case gnss.OutputData:
			metr.gnssraw = &m
		}
		metr.mutex.Unlock()
		//m.logger.Debug().Msgf("msg %+v", metr)
	}

}

func (m *Unit) publishData(metrics *metricsData) []byte {

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
	if metrics.temp != nil {
		mpb.Temperature = &pb.Metrics_Temperature{
			InBox: metrics.temp.Temperature,
		}
	}
	if metrics.internet != nil {
		mpb.Internet = &pb.Metrics_Internet{
			Connected: metrics.internet.Connected,
		}
	}

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
	m.logger.Debug().Msgf("publish %+v", mpb)

	out, err := proto.Marshal(mpb)
	if err != nil {
		m.logger.Error().Msgf("can't marshall %v", err)
	}
	return out
}
