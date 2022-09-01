package sensor

import (
	"time"

	"github.com/ci4rail/bogie-pdm/cmd/bogie-edge/internal/position"
	"github.com/ci4rail/bogie-pdm/cmd/bogie-edge/internal/steadydrive"
	pb "github.com/ci4rail/bogie-pdm/proto/go/bogie/v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type auxData struct {
	steadydrive *steadydrive.OutputData
	position    *position.OutputData
}

func (s *Unit) publisher() {
	go func() {

		s.logger.Info().Msgf("publisher running")
		var aux auxData

		triggerCh := s.ps.Sub("trigger")
		inputCh := s.ps.Sub("steadydrive", "position")

		for {
			select {
			case tr := <-triggerCh:
				// received trigger message, publish data
				o := s.publishData(&aux)
				s.logger.Debug().Msgf("received trigger %v, published %d bytes", tr, len(o))
				err := s.export.PubExport("bogie", o)
				if err != nil {
					s.logger.Error().Msgf("can't publish %v", err)
				}

			case msg := <-inputCh:
				//s.logger.Debug().Msgf("msg %v", msg)
				switch m := msg.(type) {
				case steadydrive.OutputData:
					aux.steadydrive = &m
				case position.OutputData:
					aux.position = &m
				}
			}
		}
	}()
}

func (s *Unit) publishData(aux *auxData) []byte {

	ts := time.Now()

	b := &pb.Bogie{
		Id:            int32(s.cfg.BogieID),
		TriggerType:   pb.TriggerType_UNKNOWN, // TODO
		TriggerTs:     timestamppb.New(ts),
		SensorSamples: make([]*pb.VibrationSensorSamples, len(s.sampler)),
	}
	if aux.steadydrive != nil {
		b.SteadyDrive = &pb.Bogie_SteadyDrive{
			Max: aux.steadydrive.Max[:],
			Rms: aux.steadydrive.RMS[:],
		}
	}

	if aux.position != nil {
		b.Position = &pb.Bogie_Position{
			Lat:   aux.position.Lat,
			Lon:   aux.position.Lon,
			Alt:   aux.position.Alt,
			Speed: aux.position.Speed,
		}
	}

	for samplerID, sampler := range s.sampler {

		sampler.rb.Lock()
		nSamples := sampler.rb.Buf.Len()
		samples := make([]float32, nSamples)
		timeDelta := 1 / s.cfg.SampleRate * float64(nSamples)

		for i := 0; i < nSamples; i++ {
			samples[i] = sampler.rb.Buf.At(i)
		}
		sampler.rb.Unlock()

		b.SensorSamples[samplerID] = &pb.VibrationSensorSamples{
			SensorId:      int32(samplerID),
			FirstSampleTs: timestamppb.New(ts.Add(-time.Duration(timeDelta))),
			SampleRate:    s.cfg.SampleRate,
			Samples:       samples,
		}
	}
	out, err := proto.Marshal(b)
	if err != nil {
		s.logger.Error().Msgf("can't marshall %v", err)
	}
	return out
}
