package main

import (
	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/db"
	"google.golang.org/api/option"
)

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
	const projectID = "aggregarian"

	client, ctx, err := getDatabaseClient(projectID)
	if err != nil {
		log.Fatalln(err)
		return
	}

	var data map[string]Reading
	if err := client.NewRef("log/").Get(ctx, &data); err != nil {
		log.Fatalln("Error setting value:", err)
		return
	}
	log.Printf("Data: %v", data)
}
