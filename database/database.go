package database

import (
	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/db"
	"google.golang.org/api/option"
)

func getApp(projectName string) (*firebase.App, context.Context, error) {
	ctx := context.Background()
	conf := &firebase.Config{
		DatabaseURL: fmt.Sprintf("https://%s.firebaseio.com", projectName),
	}
	fileName := fmt.Sprintf("./database/%s.json", projectName)
	opt := option.WithCredentialsFile(fileName)

	app, err := firebase.NewApp(ctx, conf, opt)
	if err != nil {
		log.Fatalln("Error initializing app:", err)
		return nil, nil, err
	}

	return app, ctx, nil
}

// GetClient firebase DB client with admin privileges
func GetClient(projectName string) (*db.Client, context.Context, error) {
	app, ctx, err := getApp(projectName)
	if err != nil {
		log.Fatalln("Error initializing database app instance:", err)
		return nil, nil, err
	}

	client, err := app.Database(ctx)
	if err != nil {
		log.Fatalln("Error initializing database client:", err)
		return nil, nil, err
	}

	return client, ctx, nil
}
