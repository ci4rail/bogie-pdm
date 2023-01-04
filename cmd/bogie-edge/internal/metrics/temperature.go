package metrics

import (
	"os"
	"strconv"
	"time"
)

const temperatureFile = "/sys/class/thermal/thermal_zone0/temp"

type temperatureData struct {
	Temperature float32
}

func (m *Unit) runTemperature() {

	for {
		time.Sleep(time.Second * 5)
		b, err := os.ReadFile(temperatureFile)
		if err != nil {
			m.logger.Error().Err(err).Msg("read temperature file")
			continue
		}
		s := string(b)
		// remove \n
		s = s[:len(s)-1]
		t, err := strconv.Atoi(s)
		if err != nil {
			m.logger.Error().Err(err).Msg("parse temperature")
			continue
		}
		m.ps.Pub(temperatureData{
			Temperature: float32(t) / 1000,
		}, "temperature")
	}
}
