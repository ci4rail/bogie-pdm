package triggerunit

import (
	"fmt"

	"github.com/cskr/pubsub"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type configuration struct {
	TriggerDuration float64 // time how long the trigger condition must be met until we trigger. In seconds
	HoldOff         float64 // time how long we wait after trigger before we trigger again. In seconds
	SteadyDrive     struct {
		Max [3]float64 // vibration must be smaller than this value
		RMS [3]float64 // vibration must be smaller than this value
	}
}

// Instance is the instance of the TriggerUnit
type Instance struct {
	cfg    *configuration
	logger zerolog.Logger
	ps     *pubsub.PubSub
}

// NewFromViper creates a new TriggerUnit from a viper configuration
func NewFromViper(viperCfg *viper.Viper, ps *pubsub.PubSub) (*Instance, error) {
	cfg, err := readConfig(viperCfg)
	if err != nil {
		return nil, err
	}
	return New(cfg, ps), nil
}

// New creates a new instance of the TriggerUnit
func New(cfg *configuration, ps *pubsub.PubSub) *Instance {

	i := &Instance{
		ps:     ps,
		cfg:    cfg,
		logger: log.With().Str("component", "triggerunit").Logger(),
	}

	i.logger.Info().Msg(fmt.Sprintf("config: %+v", cfg))

	return i
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