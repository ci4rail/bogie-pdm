package metrics

import (
	"fmt"

	"github.com/cskr/pubsub"
	"github.com/edgefarm/bogie-pdm/cmd/bogie-edge/internal/nats"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type configuration struct {
	PublishPeriod float64
}

// Unit is the instance of the MetricsUnit
type Unit struct {
	cfg    *configuration
	logger zerolog.Logger
	ps     *pubsub.PubSub
	nc     *nats.Connection
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
		ps:     ps,
		cfg:    cfg,
		logger: log.With().Str("component", "metrics").Logger(),
		nc:     nc,
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
