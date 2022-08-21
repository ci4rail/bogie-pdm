package triggerunit

import (
	"os"
	"testing"
	"time"

	"github.com/cskr/pubsub"
	"github.com/edgefarm/bogie-pdm/cmd/bogie-edge/internal/steadydrive"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func eatTrigger(ch <-chan interface{}) {
	select {
	case <-ch:
	default:
	}
}

func TestTriggerUnit(t *testing.T) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:04:05.999Z07:00"})
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	const sdTopic = "steadydrive"
	assert := assert.New(t)
	ps := pubsub.New(5)
	cfg := &configuration{
		TriggerDuration: 3,
		HoldOff:         5,
	}
	cfg.SteadyDrive.Max = [3]float64{1, 2, 3}
	cfg.SteadyDrive.RMS = [3]float64{4, 5, 7}

	tu := New(cfg, ps)

	ch := ps.Sub("trigger")

	go tu.Run()

	time.Sleep(time.Millisecond * 100)
	assert.Empty(ch, "no trigger")
	ps.Pub(steadydrive.OutputData{Timestamp: time.Now(), Max: [3]float64{4, 4, 4}, RMS: [3]float64{8, 8, 8}}, sdTopic)
	time.Sleep(time.Millisecond * 100)
	assert.Empty(ch, "no trigger")
	for i := 0; i < 3; i++ {
		assert.Empty(ch, "no trigger")
		ps.Pub(steadydrive.OutputData{Timestamp: time.Now(), Max: [3]float64{1, 1, 1}, RMS: [3]float64{1, 1, 1}}, sdTopic)
		time.Sleep(time.Millisecond * 1000)
	}
	// expect initial trigger
	time.Sleep(time.Millisecond * 100)
	assert.NotEmpty(ch, "trigger")
	eatTrigger(ch)

	// we should be in holdoff
	for i := 0; i < 8; i++ {
		assert.Empty(ch, "no trigger")
		ps.Pub(steadydrive.OutputData{Timestamp: time.Now(), Max: [3]float64{1, 1, 1}, RMS: [3]float64{1, 1, 1}}, sdTopic)
		time.Sleep(time.Millisecond * 1000)
	}
	time.Sleep(time.Millisecond * 100)
	// after holdoff, we should trigger again
	assert.NotEmpty(ch, "trigger")
	eatTrigger(ch)

}
