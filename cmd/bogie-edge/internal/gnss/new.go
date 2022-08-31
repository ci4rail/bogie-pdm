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

// New creates a new instance of SteadyDrive
func New(sub *viper.Viper, ps *pubsub.PubSub) (*Gnss, error) {
	var s Gnss
	var err error
	s.logger = log.With().Str("component", "gnss").Logger()
	s.ps = ps
	s.cfg, err = s.readConfig(sub)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (s *Gnss) readConfig(sub *viper.Viper) (*configuration, error) {
	if sub == nil {
		return nil, fmt.Errorf("missing configuration")
	}
	var cfg configuration
	err := sub.Unmarshal(&cfg)
	if err != nil {
		return nil, fmt.Errorf("unmarshal config %s", err)
	}
	s.logger.Printf("gnss config: %+v\n", cfg)

	return &cfg, nil
}
