package sender

import (
	"log"
	"time"

	"github.com/kevinnayar/aggregarian/database"
	"github.com/kevinnayar/aggregarian/sensor"
)

// Reading represents the current state of the sensor
type Reading struct {
	ReadableDate string `json:"ReadableDate,omitempty"`
	IsDry        bool
}

// Start read data from soil moisture sensor
func Start(projectName string) {
	const pin = 21
	const timezone = "America/Chicago"

	sensor, err := sensor.GetInstance(pin)
	if err != nil {
		log.Fatalln(err)
		return
	}

	client, ctx, err := database.GetClient(projectName)
	if err != nil {
		log.Fatalln(err)
		return
	}

	for {
		utc := time.Now().UTC()
		local := utc

		location, err := time.LoadLocation(timezone)
		if err != nil {
			log.Fatalln("Error setting location:", err)
			return
		}
		local = local.In(location)

		date := utc.Format("2006-01-02T15:04:05Z07:00")
		readableDate := local.Format("Mon, 02 Jan 2006 15:04:05 MST")

		isDry := true
		if sensor.Read() == 0 {
			isDry = false
		}

		reference := "log/" + date

		if err := client.NewRef(reference).Set(ctx, &Reading{
			ReadableDate: readableDate,
			IsDry:        isDry,
		}); err != nil {
			log.Fatalln("Error setting value:", err)
			return
		}

		log.Printf("Date: %v, ReadableDate: %v, IsDry: %v", date, readableDate, isDry)

		time.Sleep(6 * time.Hour)
	}
}
