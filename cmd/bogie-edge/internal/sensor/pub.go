package sensor

import (
	"time"

	"github.com/edgefarm/bogie-pdm/cmd/bogie-edge/internal/position"
	"github.com/edgefarm/bogie-pdm/cmd/bogie-edge/internal/steadydrive"
	pb "github.com/edgefarm/bogie-pdm/proto/go/bogie/v1"
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
				s.logger.Debug().Msgf("received trigger %v", tr)
				s.publish(&aux)

			case msg := <-inputCh:
				s.logger.Debug().Msgf("msg %v", msg)
				switch m := msg.(type) {
				case steadydrive.OutputData:
					aux.steadydrive = &m
				}
			}
		}
	}()
}

func (s *Unit) publish(aux *auxData) {

	ts := time.Now()

	b := &pb.Bogie{
		Id:            1,                      // TODO
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
	// if aux.position != nil {
	// 	// TODO
	// }

	for samplerID, sampler := range s.sampler {

		nSamples := sampler.rb.Buf.Len()
		samples := make([]float32, nSamples)
		timeDelta := 1 / s.cfg.SampleRate * float64(nSamples)

		sampler.rb.Lock()
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
	_ = out
}
