package sensor

import (
	"log"

	"github.com/stianeikeland/go-rpio"
)

// GetInstance instance of soil moisture sensor
func GetInstance(pin int) (rpio.Pin, error) {
	if err := rpio.Open(); err != nil {
		log.Fatal("Error initializing soil moisture sensor:", err)
		return rpio.Pin(pin), err
	}

	sensor := rpio.Pin(pin)
	sensor.Input()
	sensor.PullUp()
	sensor.Detect(rpio.FallEdge)

	return sensor, nil
}
