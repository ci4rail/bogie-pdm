package triggerunit

import (
	"time"

	"github.com/edgefarm/bogie-pdm/cmd/bogie-edge/internal/steadydrive"
)

// OutputData is the output data of the trigger unit.
// It is published when triggered
type OutputData struct {
	// empty
}

type fsmState int

const (
	wait    fsmState = iota // waiting for trigger condition
	arm                     // trigger condition met, must hold for n seconds
	holdoff                 // wait until holdoff time over
)

// Run starts the trigger unit.
// should be called from a goroutine.
func (i *Instance) Run() {
	i.logger.Info().Msg("running")

	inputCh := i.ps.Sub("steadydrive")

	state := wait
	prevState := state
	stateEntryTime := time.Now()
	var curSteadydrive *steadydrive.OutputData

	for {
		select {
		case <-time.After(time.Millisecond * 500):
			i.logger.Debug().Msg("timer")

		case msg := <-inputCh:
			i.logger.Debug().Msgf("msg %v", msg)
			switch m := msg.(type) {
			case steadydrive.OutputData:
				curSteadydrive = &m
			}
		}

		switch state {
		case wait:
			if i.isTriggerMet(curSteadydrive) {
				state = arm
				i.logger.Debug().Msg("trigger met")
			}
		case arm:
			if !i.isTriggerMet(curSteadydrive) {
				state = wait
				i.logger.Debug().Msg("trigger lost")
			}

			if time.Since(stateEntryTime) > time.Duration(i.cfg.TriggerDuration)*time.Second {
				i.ps.Pub(OutputData{}, "trigger")
				state = holdoff
			}
		case holdoff:
			if time.Since(stateEntryTime) > time.Duration(i.cfg.HoldOff)*time.Second {
				state = wait
			}
		}
		if state != prevState {
			i.logger.Debug().Msgf("state changed from %d to %d", prevState, state)
			stateEntryTime = time.Now()
		}

		prevState = state
	}
}

func (i *Instance) isTriggerMet(sd *steadydrive.OutputData) bool {
	return i.isSteadyDriveOk(sd)
}

func (i *Instance) isSteadyDriveOk(sd *steadydrive.OutputData) bool {
	if sd == nil {
		return false
	}
	for ax := 0; ax < 3; ax++ {
		if sd.Max[ax] > i.cfg.SteadyDrive.Max[ax] || sd.RMS[ax] > i.cfg.SteadyDrive.RMS[ax] {
			return false
		}
	}
	return true
}
