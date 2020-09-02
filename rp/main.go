package main

import (
	"context"
	"fmt"
	"log"
	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/db"
	"github.com/stianeikeland/go-rpio"
	"google.golang.org/api/option"
)

func getSensorInstance(pin int) rpio.Pin {
	if err := rpio.Open(); err != nil {
		log.Fatal("Error initializing soil moisture sensor:", err)
	}

	sensor := rpio.Pin(pin)
	sensor.Input()
	sensor.PullUp()
	sensor.Detect(rpio.FallEdge)

	return sensor
}

func getDatabaseClient(projectName string) (*db.Client, context.Context) {
	ctx := context.Background()
	conf := &firebase.Config{
		DatabaseURL: fmt.Sprintf("https://%s.firebaseio.com", projectName),
	}
	fileName := fmt.Sprintf("%s.json", projectName)
	opt := option.WithCredentialsFile(fileName)

	app, err := firebase.NewApp(ctx, conf, opt)
	if err != nil {
		log.Fatalln("Error initializing app:", err)
	}

	client, err := app.Database(ctx)
	if err != nil {
		log.Fatalln("Error initializing database client:", err)
	}

	return client, ctx
}

// Reading represents the current state of the sensor
type Reading struct {
	ReadableDate string `json:"ReadableDate,omitempty"`
	IsDry        bool
}

func main() {
	PIN := 21
	PROJECTNAME := "aggregarian"
	TIMEZONE := "America/Chicago"

	sensor := getSensorInstance(PIN)
	client, ctx := getDatabaseClient(PROJECTNAME)

	for {
		utc := time.Now().UTC()
		local := utc

		location, locationErr := time.LoadLocation(TIMEZONE)
		if locationErr != nil {
			log.Fatalln("Error setting location:", locationErr)
		} else {
			local = local.In(location)
		}

		date := utc.Format("2006-01-02T15:04:05Z07:00")
		readableDate := local.Format("Mon, 02 Jan 2006 15:04:05 MST")

		isDry := true
		if sensor.Read() == 0 {
			isDry = false
		}

		reference := "log/" + date
		referenceErr := client.NewRef(reference).Set(ctx, &Reading{
			ReadableDate: readableDate,
			IsDry:        isDry,
		})

		if referenceErr != nil {
			log.Fatalln("Error setting value:", referenceErr)
		}

		log.Printf("Date: %v, ReadableDate: %v, IsDry: %v", date, readableDate, isDry)

		time.Sleep(6 * time.Hour)
	}
}
