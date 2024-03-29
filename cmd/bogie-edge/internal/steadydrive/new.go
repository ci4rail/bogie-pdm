package steadydrive

import (
	"fmt"
	"time"

	"github.com/ci4rail/io4edge-client-go/motionsensor"
	"github.com/cskr/pubsub"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type configuration struct {
	OutputDataRateHz float64
	MotionSensor     struct {
		DeviceAddress  string
		SampleRate     float64
		FullScale      int32
		HighPassFilter bool
		BandwidthRatio int32
	}
}

// SteadyDrive is the instance of the SteadyDrive
type SteadyDrive struct {
	cfg           *configuration
	io4edgeClient *motionsensor.Client
	logger        zerolog.Logger
	ps            *pubsub.PubSub
}

// New creates a new instance of SteadyDrive
func New(sub *viper.Viper, ps *pubsub.PubSub) (*SteadyDrive, error) {
	var s SteadyDrive
	var err error
	s.logger = log.With().Str("component", "steadydrive").Logger()
	s.ps = ps
	s.cfg, err = s.readConfig(sub)
	if err != nil {
		return nil, err
	}
	if err := s.configMotionSensor(); err != nil {
		return nil, err
	}
	return &s, nil
}

func (s *SteadyDrive) readConfig(sub *viper.Viper) (*configuration, error) {
	if sub == nil {
		return nil, fmt.Errorf("missing configuration")
	}
	var cfg configuration
	err := sub.Unmarshal(&cfg)
	if err != nil {
		return nil, fmt.Errorf("unmarshal config %s", err)
	}
	s.logger.Printf("steadydrive config: %+v\n", cfg)

	return &cfg, nil
}

func (s *SteadyDrive) configMotionSensor() error {
	timeout := time.Duration(0)
	c, err := motionsensor.NewClientFromUniversalAddress(s.cfg.MotionSensor.DeviceAddress, timeout)
	if err != nil {
		return fmt.Errorf("failed to create motionsensor client: %v", err)
	}
	s.io4edgeClient = c

	// set configuration
	if err := c.UploadConfiguration(
		motionsensor.WithSampleRate(uint32(s.cfg.MotionSensor.SampleRate*1000.0)),
		motionsensor.WithFullScale(s.cfg.MotionSensor.FullScale),
		motionsensor.WithHighPassFilterEnable(s.cfg.MotionSensor.HighPassFilter),
		motionsensor.WithBandWidthRatio(s.cfg.MotionSensor.BandwidthRatio)); err != nil {
		return fmt.Errorf("failed to set motionsensor configuration: %v", err)
	}

	return nil
}
