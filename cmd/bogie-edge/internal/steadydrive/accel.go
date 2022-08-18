package steadydrive

import (
	"fmt"
	"log"
	"time"

	"github.com/ci4rail/io4edge-client-go/functionblock"
	mspb "github.com/ci4rail/io4edge_api/motionSensor/go/motionSensor/v1"
	"github.com/edgefarm/bogie-pdm/pkg/signalprocessing"
)

// OutputData is the output data of the SteadyDrive
type OutputData struct {
	Max [3]float64 // x, y, z
	RMS [3]float64 // x, y, z
}

func outputDataFromAccelerometerValues(samples []*mspb.Sample) *OutputData {
	var outputData OutputData

	for axis := 0; axis < 3; axis++ {
		var data = make([]float64, len(samples))

		for index, sample := range samples {
			switch axis {
			case 0:
				data[index] = float64(sample.X)
			case 1:
				data[index] = float64(sample.Y)
			case 2:
				data[index] = float64(sample.Z)
			}
		}
		outputData.Max[axis] = signalprocessing.Max(data)
		outputData.RMS[axis] = signalprocessing.RMS(data)
	}
	return &outputData
}

// Run runs the SteadyDrive function block
func (i *Instance) Run() error {

	c := i.io4edgeClient
	// start stream
	err := c.StartStream(
		functionblock.WithBucketSamples(10),
		functionblock.WithBufferedSamples(200),
	)
	if err != nil {
		return fmt.Errorf("startStream failed: %v", err)
	}

	for {
		samples := make([]*mspb.Sample, 0)
		start := time.Now()

		for {
			sd, err := c.ReadStream(time.Second * 2)
			if err != nil {
				log.Printf("%s: readStream failed: %v", tag, err)
				continue
			}
			samples = append(samples, sd.FSData.GetSamples()...)
			if time.Since(start) > time.Second/time.Duration(i.cfg.OutputDataRateHz) {
				break
			}
		}
		o := outputDataFromAccelerometerValues(samples)
		log.Printf("%s: %d samples, %+v\n", tag, len(samples), *o)
	}
}
