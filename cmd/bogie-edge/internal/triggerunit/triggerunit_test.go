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

func TestTriggerUnit(t *testing.T) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
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
	start := time.Now()
	// simulate steadydrive
	go func() {
		time.Sleep(time.Second * 1)
		ps.Pub(steadydrive.OutputData{Max: [3]float64{4, 4, 4}, RMS: [3]float64{8, 8, 8}}, sdTopic)
		time.Sleep(time.Second * 1)
		ps.Pub(steadydrive.OutputData{Max: [3]float64{4, 4, 4}, RMS: [3]float64{8, 8, 8}}, sdTopic)
		time.Sleep(time.Second * 1)
		ps.Pub(steadydrive.OutputData{Max: [3]float64{1, 1, 1}, RMS: [3]float64{1, 1, 1}}, sdTopic)
		time.Sleep(time.Second * 5)
		ps.Pub(steadydrive.OutputData{Max: [3]float64{1, 1, 1}, RMS: [3]float64{1, 1, 1}}, sdTopic)
		time.Sleep(time.Second * 5)
		ps.Pub(steadydrive.OutputData{Max: [3]float64{1, 1, 1}, RMS: [3]float64{1, 1, 1}}, sdTopic)
		time.Sleep(time.Second * 5)
		ps.Pub(steadydrive.OutputData{Max: [3]float64{4, 4, 4}, RMS: [3]float64{8, 8, 8}}, sdTopic)
	}()

	<-ch
	assert.InDelta(4.0, time.Since(start), 0.001)
}
