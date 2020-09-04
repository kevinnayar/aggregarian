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

func getSensorInstance(pin int) (rpio.Pin, error) {
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

func getDatabaseClient(projectName string) (*db.Client, context.Context, error) {
	ctx := context.Background()
	conf := &firebase.Config{
		DatabaseURL: fmt.Sprintf("https://%s.firebaseio.com", projectName),
	}
	fileName := fmt.Sprintf("../config/%s.json", projectName)
	opt := option.WithCredentialsFile(fileName)

	app, err := firebase.NewApp(ctx, conf, opt)
	if err != nil {
		log.Fatalln("Error initializing app:", err)
		return nil, nil, err
	}

	client, err := app.Database(ctx)
	if err != nil {
		log.Fatalln("Error initializing database client:", err)
		return nil, nil, err
	}

	return client, ctx, nil
}

// Reading represents the current state of the sensor
type Reading struct {
	ReadableDate string `json:"ReadableDate,omitempty"`
	IsDry        bool
}

func main() {
	const pin = 21
	const projectID = "aggregarian"
	const timezone = "America/Chicago"

	sensor, err := getSensorInstance(pin)
	if err != nil {
		log.Fatalln(err)
		return
	}

	client, ctx, err := getDatabaseClient(projectID)
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
