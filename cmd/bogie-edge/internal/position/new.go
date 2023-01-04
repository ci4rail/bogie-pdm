package position

import (
	"fmt"

	"github.com/cskr/pubsub"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type configuration struct {
	MaxGnssAge float64 // if valid gnns data is older than this, report no position (seconds)
}

// Unit is the instance of the Unit
type Unit struct {
	cfg    *configuration
	logger zerolog.Logger
	ps     *pubsub.PubSub
}

// NewFromViper creates a new PositionUnit from a viper configuration
func NewFromViper(viperCfg *viper.Viper, ps *pubsub.PubSub) (*Unit, error) {
	cfg, err := readConfig(viperCfg)
	if err != nil {
		return nil, err
	}
	return New(cfg, ps), nil
}

// New creates a new instance of PositionUnit
func New(cfg *configuration, ps *pubsub.PubSub) *Unit {

	p := &Unit{
		ps:     ps,
		cfg:    cfg,
		logger: log.With().Str("component", "position").Logger(),
	}

	p.logger.Info().Msg(fmt.Sprintf("config: %+v", cfg))
	return p
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
