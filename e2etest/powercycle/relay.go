package powercycle

import (
	"github.com/Tinkerforge/go-api-bindings/industrial_dual_relay_bricklet"
	"github.com/Tinkerforge/go-api-bindings/ipconnection"
)

const addr = "192.168.24.14:4223"
const uid = "NuQ"
const channel = 0

func setRelay(state bool) error {
	ipcon := ipconnection.New()
	defer ipcon.Close()
	idr, err := industrial_dual_relay_bricklet.New(uid, &ipcon) // Create device object.
	if err != nil {
		return err
	}

	ipcon.Connect(addr) // Connect to brickd.
	defer ipcon.Disconnect()

	idr.SetSelectedValue(channel, state)
	return nil
}
