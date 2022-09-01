package sensor

import (
	"fmt"

	"github.com/ci4rail/bogie-pdm/cmd/bogie-edge/internal/nats"
	"github.com/cskr/pubsub"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type configuration struct {
	BogieID        int
	DeviceAddress  []string
	SampleRate     float64
	RingBufEntries int32
}

type sampler struct {
	deviceAddress string
	rb            *samplesRingbuf
}

// Unit is the instance of the SensorUnit
type Unit struct {
	cfg     *configuration
	logger  zerolog.Logger
	ps      *pubsub.PubSub
	sampler []*sampler
	nc      *nats.Connection
}

// NewFromViper creates a new MetricsUnit from a viper configuration
func NewFromViper(viperCfg *viper.Viper, ps *pubsub.PubSub, nc *nats.Connection) (*Unit, error) {
	cfg, err := readConfig(viperCfg)
	if err != nil {
		return nil, err
	}
	return New(cfg, ps, nc), nil
}

// New creates a new instance of the MetricsUnit
func New(cfg *configuration, ps *pubsub.PubSub, nc *nats.Connection) *Unit {

	t := &Unit{
		ps:      ps,
		cfg:     cfg,
		logger:  log.With().Str("component", "sensor").Logger(),
		sampler: make([]*sampler, len(cfg.DeviceAddress)),
		nc:      nc,
	}

	t.logger.Info().Msg(fmt.Sprintf("config: %+v", cfg))

	return t
}

func readConfig(sub *viper.Viper) (*configuration, error) {
	if sub == nil {
		return nil, fmt.Errorf("missing configuration")
	}
	var cfg configuration
	err := sub.Unmarshal(&cfg)
	if err != nil {
		return nil, fmt.Errorf("unmarshal config %s", err)
	}

	return &cfg, nil
}

// Run starts the sensor unit
// It returns as soon as all go routines are started
func (s *Unit) Run() error {
	s.logger.Info().Msg("sensorunit starting")

	// start sensor sampling
	for i, addr := range s.cfg.DeviceAddress {
		rb, err := s.sample(addr, s.cfg.SampleRate, s.cfg.RingBufEntries)
		if err != nil {
			return err
		}

		s.sampler[i] = &sampler{
			deviceAddress: addr,
			rb:            rb,
		}
	}
	// start publisher
	s.logger.Info().Msg("about to start publisher")
	s.publisher()

	return nil
}
