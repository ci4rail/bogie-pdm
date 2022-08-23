package sensor

import (
	"fmt"

	"github.com/cskr/pubsub"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type configuration struct {
	DeviceAddress  []string
	SampleRate     float64
	RingBufEntries int32
}

// Unit is the instance of the SensorUnit
type Unit struct {
	cfg    *configuration
	logger zerolog.Logger
	ps     *pubsub.PubSub
}

// NewFromViper creates a new MetricsUnit from a viper configuration
func NewFromViper(viperCfg *viper.Viper, ps *pubsub.PubSub) (*Unit, error) {
	cfg, err := readConfig(viperCfg)
	if err != nil {
		return nil, err
	}
	return New(cfg, ps), nil
}

// New creates a new instance of the MetricsUnit
func New(cfg *configuration, ps *pubsub.PubSub) *Unit {

	t := &Unit{
		ps:     ps,
		cfg:    cfg,
		logger: log.With().Str("component", "sensor").Logger(),
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
