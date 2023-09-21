package metrics

import (
	"time"

	"github.com/maltegrosse/go-modemmanager"
)

type modemData struct {
	OperatorName  string
	SignalQuality float32
	CellId        string
}

func (m *Unit) runModem() {
	mmgr, err := modemmanager.NewModemManager()
	if err != nil {
		m.logger.Error().Err(err).Msg("could not create modem manager")
		return
	}

	for {
		time.Sleep(time.Second * 5)
		modems, err := mmgr.GetModems()
		if err != nil {
			m.logger.Error().Err(err).Msg("get modems")
			continue
		}
		for _, modem := range modems {
			modem3gpp, err := modem.Get3gpp()
			if err != nil {
				m.logger.Error().Err(err).Msg("get 3gpp")
				continue
			}
			opName, err := modem3gpp.GetOperatorName()
			if err != nil {
				m.logger.Error().Err(err).Msg("get operator name")
				continue
			}

			location, err := modem.GetLocation()
			if err != nil {
				m.logger.Error().Err(err).Msg("get location")
				continue
			}

			currentLocation, err := location.GetCurrentLocation()
			if err != nil {
				m.logger.Error().Err(err).Msg("get current location")
				continue
			}
			Ci := currentLocation.ThreeGppLacCi.Ci

			signalQuality, _, err := modem.GetSignalQuality()
			if err != nil {
				m.logger.Error().Err(err).Msg("get signal quality")
				continue
			}
			m.ps.Pub(modemData{
				OperatorName:  opName,
				SignalQuality: float32(signalQuality),
				CellId:        Ci,
			}, "cellular")
			break // only one modem for now
		}
	}
}
