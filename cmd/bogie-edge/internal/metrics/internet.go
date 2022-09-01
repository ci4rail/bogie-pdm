package metrics

import (
	"net"
	"time"
)

const lookupURI = "www.wikipedia.org"

type internetData struct {
	Connected bool
}

func (m *Unit) runInternet() {

	data := internetData{}
	for {
		time.Sleep(time.Second * 5)
		_, err := net.LookupIP(lookupURI)
		data.Connected = err == nil
		m.ps.Pub(data, "internet")
	}
}
