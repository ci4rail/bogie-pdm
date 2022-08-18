package steadydrive

import (
	"fmt"

	"github.com/spf13/viper"
)

type configuration struct {
	accelDeviceAddress string
	outputDataRateHz   float64
}

// Instance is the instance of the SteadyDrive
type Instance struct {
	cfg *configuration
}

// New creates a new instance of SteadyDrive
func New(sub *viper.Viper) (*Instance, error) {
	var i Instance
	var err error
	i.cfg, err = readConfig(sub)
	if err != nil {
		return nil, err
	}
	return &i, nil
}

func readConfig(sub *viper.Viper) (*configuration, error) {
	if sub == nil {
		return nil, fmt.Errorf("missing configuration")
	}
	var cfg configuration
	err := sub.Unmarshal(&cfg)
	if err != nil {
		return nil, fmt.Errorf("unmarhsal config %s", err)
	}
	return &cfg, nil
}
