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
		CompareType int // 0 = vibration must be smaller than Max/RMS, 1 = vibration must be larger than Max/RMS
		Max         [3]float64
		RMS         [3]float64
	}
	Position struct {
		MinLat   float64 // minimum latitude
		MaxLat   float64 // maximum latitude
		MinLon   float64 // minimum longitude
		MaxLon   float64 // maximum longitude
		MinSpeed float64 // minimum speed
		MaxSpeed float64 // maximum speed
	}
}

// TriggerUnit is the instance of the TriggerUnit
type TriggerUnit struct {
	cfg    *configuration
	logger zerolog.Logger
	ps     *pubsub.PubSub
}

// NewFromViper creates a new TriggerUnit from a viper configuration
func NewFromViper(viperCfg *viper.Viper, ps *pubsub.PubSub) (*TriggerUnit, error) {
	cfg, err := readConfig(viperCfg)
	if err != nil {
		return nil, err
	}
	return New(cfg, ps), nil
}

// New creates a new instance of the TriggerUnit
func New(cfg *configuration, ps *pubsub.PubSub) *TriggerUnit {

	t := &TriggerUnit{
		ps:     ps,
		cfg:    cfg,
		logger: log.With().Str("component", "triggerunit").Logger(),
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
