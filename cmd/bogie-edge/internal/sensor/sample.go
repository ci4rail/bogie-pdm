package sensor

import (
	"fmt"
	"sync"
	"time"

	anain "github.com/ci4rail/io4edge-client-go/analogintypea"
	"github.com/ci4rail/io4edge-client-go/functionblock"
	"github.com/edgefarm/bogie-pdm/pkg/ringbuf"
)

type samplesRingbuf struct {
	Buf   *ringbuf.Ringbuf[float32]
	mutex *sync.Mutex
}

func (s *samplesRingbuf) Lock() {
	s.mutex.Lock()
}
func (s *samplesRingbuf) Unlock() {
	s.mutex.Unlock()
}

func (s *Unit) sample(deviceAddress string, sampleRate float64, rbEntries int32) (*samplesRingbuf, error) {
	s.logger.Info().Msgf("%s: sampler starting", deviceAddress)

	// create ringbuf
	rb := &samplesRingbuf{
		Buf:   ringbuf.New[float32](int(rbEntries)),
		mutex: &sync.Mutex{},
	}

	timeout := time.Duration(0)
	c, err := anain.NewClientFromUniversalAddress(deviceAddress, timeout)
	if err != nil {
		return nil, fmt.Errorf("failed to create anain %s client: %v", deviceAddress, err)
	}
	// set configuration
	if err := c.UploadConfiguration(
		anain.WithSampleRate(uint32(sampleRate)),
	); err != nil {
		return nil, fmt.Errorf("failed to set anain %s configuration: %v", deviceAddress, err)
	}

	// sample in background
	go func() {
		// start stream
		err := c.StartStream(
			functionblock.WithBucketSamples(100),
			functionblock.WithBufferedSamples(200),
		)
		if err != nil {
			s.logger.Error().Msgf("startStream failed for %s: %v", deviceAddress, err)
		}

		for {

			for {
				sd, err := c.ReadStream(time.Second * 2)
				if err != nil {
					s.logger.Error().Err(err).Msgf("readstream %d failed", deviceAddress)
					continue
				}
				samples := sd.FSData.GetSamples()
				//s.logger.Debug().Msgf("%s: read %d samples", deviceAddress, len(samples))

				tsLast := uint64(0)
				periodUS := 1e6 / uint64(sampleRate)

				// copy samples to ringbuf
				rb.Lock()
				for i, sample := range samples {
					if i != 0 {
						if sample.Timestamp-tsLast > periodUS*2 {
							s.logger.Warn().Msgf("timestamp gap detected for %s: %d", deviceAddress, sample.Timestamp-tsLast)
						}
					}
					rb.Buf.Push(float32(sample.Value))
					tsLast = sample.Timestamp
				}
				rb.Unlock()
			}
		}
	}()
	return rb, nil
}
