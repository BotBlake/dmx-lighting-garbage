package serialhw

import "go.bug.st/serial"

type Device struct {
	Name string
}

func ListDevices() ([]Device, error) {
	ports, err := serial.GetPortsList()
	if err != nil {
		return nil, err
	}

	devs := make([]Device, 0, len(ports))
	for _, p := range ports {
		devs = append(devs, Device{Name: p})
	}
	return devs, nil
}
