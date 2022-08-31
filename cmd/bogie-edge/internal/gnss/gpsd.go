package gnss

import (
	"time"

	"github.com/edgefarm/edgefarm-service-modules/cmd/location-provider/pkg/gpsd"
)

// OutputData is the output data of the Gnss function block
type OutputData struct {
	Timestamp time.Time
	Lat       float64
	Lon       float64
	Alt       float64 // altitude in meters
	Speed     float64 // in m/s
	Eph       float64 // Estimated horizontal Position (2D) Error in meters
	Mode      int     // 0=unknown, 1=no fix, 2=2D, 3=3D.
	NumSats   int     // Number of satellites used in the solution
}

// Run runs the Gnss function block
// no need to run it in a separate goroutine
func (g *Gnss) Run() error {
	const connectTimeoutSeconds = 10
	gpsChan := make(chan *gpsd.Connection)
	go func() {
		gpsdHost := g.cfg.GpsdURL
		for i := 0; i < connectTimeoutSeconds; i++ {
			if gpsdClient, err := gpsd.NewClient(gpsdHost); err != nil {
				g.logger.Warn().Msg("Can't connect to gpsd server, reconnecting")
			} else {
				g.logger.Info().Msg("connected to gpsd server")
				gpsChan <- gpsdClient
				return
			}
			time.Sleep(time.Second)
		}
		g.logger.Error().Msg("Can't connect to gpsd server")
	}()

	gpsClient := <-gpsChan
	var o OutputData

	// send output data whenever gpsd updates TPV data
	gpsClient.RegisterTpv(func(r interface{}) {
		tpv := r.(*gpsd.Tpv)

		o.Timestamp = time.Now()
		o.Lat = tpv.Lat
		o.Lon = tpv.Lon
		o.Alt = tpv.Althae
		o.Speed = tpv.Speed
		o.Eph = tpv.Eph
		o.Mode = tpv.Mode

		g.ps.Pub(o, "gnssraw")
	})

	// when sky data is available, update number of satellites in output data, but don't send it
	gpsClient.RegisterSky(func(r interface{}) {
		sky := r.(*gpsd.Sky)
		o.NumSats = numSatellites(sky)
	})

	_, err := gpsClient.Watch()
	return err
}

func numSatellites(sky *gpsd.Sky) int {
	var num int
	for _, s := range sky.Satellites {
		if s.Used {
			num++
		}
	}
	return num
}
