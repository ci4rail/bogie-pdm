package triggerunit

import (
	"time"

	"github.com/edgefarm/bogie-pdm/cmd/bogie-edge/internal/position"
	"github.com/edgefarm/bogie-pdm/cmd/bogie-edge/internal/steadydrive"
)

// OutputData is the output data of the trigger unit.
// It is published when triggered
type OutputData struct {
	TriggerType string // "triggered" or "manual"
}

type fsmState int
type fsmEvent int

const (
	wait    fsmState = iota // waiting for trigger condition
	arm                     // trigger condition met, must hold for n seconds
	holdoff                 // wait until holdoff time over
)

const (
	none fsmEvent = iota
	timer
	message
)

const hugeTime = time.Duration(1<<63 - 1)

type inputData struct {
	steadydrive *steadydrive.OutputData
	position    *position.OutputData
}

// Run starts the trigger unit.
// should be called from a goroutine.
func (t *TriggerUnit) Run() {
	t.logger.Info().Msg("running")

	inputCh := t.ps.Sub("steadydrive", "position")

	state := wait
	prevState := state
	var inputs inputData
	stateTimer := time.Duration(hugeTime)
	stateTimerStart := time.Now()
	stateEntryTime := time.Now()
	initCh := make(chan struct{}, 1)
	initCh <- struct{}{}

	for {
		event := none
		select {
		case <-time.After(stateTimer):
			//t.logger.Debug().Msg("timer")
			event = timer

		case msg := <-inputCh:
			//t.logger.Debug().Msgf("msg %v", msg)
			switch m := msg.(type) {
			case steadydrive.OutputData:
				inputs.steadydrive = &m
			case position.OutputData:
				inputs.position = &m
			}
			event = message
		case <-initCh:
			//t.logger.Debug().Msg("initCh")
			event = message
		}

		switch state {
		case wait:
			if event == message && t.isTriggerMet(&inputs) {
				state = arm
				t.logger.Debug().Msg("trigger met")
			}
		case arm:
			if !t.isTriggerMet(&inputs) {
				state = wait
				t.logger.Debug().Msg("trigger lost")
			} else if event == timer {
				if time.Since(stateEntryTime) > time.Duration(t.cfg.TriggerDuration)*time.Second {
					t.ps.Pub(OutputData{TriggerType: "triggered"}, "trigger")
					state = holdoff
					t.logger.Info().Msg("trigger")
				}
			}
		case holdoff:
			if event == timer {
				state = wait
			}
		}
		if state != prevState {
			t.logger.Debug().Msgf("state changed from %d to %d", prevState, state)

			switch state {
			case wait:
				stateTimer = time.Duration(hugeTime)
			case arm:
				stateTimer = time.Second
			case holdoff:
				stateTimer = time.Duration(t.cfg.HoldOff) * time.Second
			}
			stateEntryTime = time.Now()
			stateTimerStart = time.Now()
			initCh <- struct{}{}
		}
		if event != timer {
			if state != arm {
				stateTimer -= time.Since(stateTimerStart)
				if stateTimer < 0 {
					stateTimer = 0
				}
				stateTimerStart = time.Now()
			}
		}
		prevState = state
	}
}

func (t *TriggerUnit) isTriggerMet(inputs *inputData) bool {
	return t.isSteadyDriveOk(inputs.steadydrive) && t.isPositionOk(inputs.position)
}

func (t *TriggerUnit) isSteadyDriveOk(sd *steadydrive.OutputData) bool {
	if sd == nil {
		t.logger.Debug().Msg("steadydrive nil")
		return false
	}
	if time.Since(sd.Timestamp) > time.Second*2 {
		t.logger.Debug().Msg("steadydrive old")
		return false // ignore old data
	}

	for ax := 0; ax < 3; ax++ {
		if sd.Max[ax] > t.cfg.SteadyDrive.Max[ax] || sd.RMS[ax] > t.cfg.SteadyDrive.RMS[ax] {
			t.logger.Debug().Msgf("max/rms ax:%d", ax)
			return false
		}
	}
	return true
}

func (t *TriggerUnit) isPositionOk(p *position.OutputData) bool {
	if p == nil {
		return false
	}
	if time.Since(p.Timestamp) > time.Second*2 {
		return false // ignore old data
	}
	if !p.Valid {
		return false
	}
	if p.Lat < t.cfg.Position.MinLat || p.Lat > t.cfg.Position.MaxLat ||
		p.Lon < t.cfg.Position.MinLon || p.Lon > t.cfg.Position.MaxLon {
		return false
	}
	if p.Speed < t.cfg.Position.MinSpeed || p.Speed > t.cfg.Position.MaxSpeed {
		return false
	}

	return true
}
