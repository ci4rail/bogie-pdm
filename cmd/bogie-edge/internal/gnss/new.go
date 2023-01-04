package gnss

import (
	"fmt"

	"github.com/cskr/pubsub"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type configuration struct {
	GpsdURL string // e.g. "localhost:2947"
}

// Gnss is the instance of the Gnss
type Gnss struct {
	cfg    *configuration
	logger zerolog.Logger
	ps     *pubsub.PubSub
}

// NewFromViper creates a new Gnss Unit from a viper configuration
func NewFromViper(viperCfg *viper.Viper, ps *pubsub.PubSub) (*Gnss, error) {
	cfg, err := readConfig(viperCfg)
	if err != nil {
		return nil, err
	}
	return New(cfg, ps), nil
}

// New creates a new instance of Gnss Unit
func New(cfg *configuration, ps *pubsub.PubSub) *Gnss {

	g := &Gnss{
		ps:     ps,
		cfg:    cfg,
		logger: log.With().Str("component", "gnss").Logger(),
	}

	g.logger.Info().Msg(fmt.Sprintf("config: %+v", cfg))
	return g
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
