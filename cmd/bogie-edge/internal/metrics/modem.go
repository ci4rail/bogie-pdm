package metrics

import (
	"fmt"
	"log"
	"time"

	"github.com/maltegrosse/go-modemmanager"
)

type modemData struct {
	OperatorName  string
	SignalQuality float32
}

func (m *Unit) runModem() {
	mmgr, err := modemmanager.NewModemManager()
	if err != nil {
		log.Fatal(err.Error())
	}

	for {
		modems, err := mmgr.GetModems()
		if err != nil {
			log.Fatal(err.Error())
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
			fmt.Println(" - OperatorName: ", opName)

			signalQuality, recent, err := modem.GetSignalQuality()
			if err != nil {
				m.logger.Error().Err(err).Msg("get signal quality")
				continue
			}
			fmt.Println(" - SignalQuality: ", signalQuality, recent)
			m.ps.Pub(&modemData{
				OperatorName:  opName,
				SignalQuality: float32(signalQuality),
			})
			break // only one modem for now
		}
		time.Sleep(time.Second * 5)
	}
}
