package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/db"
	"github.com/stianeikeland/go-rpio"
	"google.golang.org/api/option"
)

func initSensor(pin int) rpio.Pin {
	if err := rpio.Open(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	sensor := rpio.Pin(pin)
	sensor.Input()
	sensor.PullUp()
	sensor.Detect(rpio.FallEdge)

	return sensor
}

func initFirebase(projectID string) (*db.Client, context.Context) {
	ctx := context.Background()
	conf := &firebase.Config{
		DatabaseURL: fmt.Sprintf("https://%s.firebaseio.com", projectID),
	}
	fileName := fmt.Sprintf("%s.json", projectID)
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

func main() {
	PIN := 21
	sensor := initSensor(PIN)
	client, ctx := initFirebase("aggregarian")

	for {
		t := time.Now()
		formattedTime := t.Format("2006-01-02_15:04:05")

		isDry := true
		if sensor.Read() == 0 {
			isDry = false
		}

		ref := client.NewRef(formattedTime)
		err := ref.Set(ctx, isDry)

		if err != nil {
			log.Fatalln("Error setting value:", err)
		}

		log.Printf("Dry = %v, Timestamp = %v", isDry, formattedTime)

		time.Sleep(6 * time.Hour)
	}
}
