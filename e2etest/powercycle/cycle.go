package powercycle

import (
	"fmt"
	"math/rand"
	"time"

	fw "github.com/edgefarm/edgefarm.core/test/pkg/framework"
	g "github.com/onsi/ginkgo/v2"
)

var _ = g.Describe("PowerCycle", g.Serial, func() {

	// var (
	// 	f *fw.Framework
	// )
	// g.BeforeEach(func() {
	// 	f = fw.DefaultFramework
	// })
	// g.AfterEach(func() {
	// })

	g.It("App publishes data to NATS after powercycle", func() {
		fw.ExpectNoError(setRelay(false)) // turn off device
		time.Sleep(2 * time.Second)
		fw.ExpectNoError(setRelay(true)) // turn on device
		ts := time.Now()
		natsSub, err := NewNatsSubscriber(natsURL, credsFile, exportSubject, "e2ecycleC", streamName)
		fw.ExpectNoError(err)
		defer natsSub.Close()

		for {
			ready, err := natsSub.WaitMsg(10*time.Second, ts)
			fw.ExpectNoError(err)

			if ready {
				break
			}
			if time.Since(ts) > 15*time.Minute {
				g.Fail("No valid message received after 15 minutes")
			}
		}
		// Wait random time
		fmt.Print("Wait random time")
		time.Sleep(time.Duration(rand.Intn(60)) * time.Second)
	})
})
