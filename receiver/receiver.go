package receiver

import (
	"log"

	"github.com/kevinnayar/aggregarian/database"
)

// Reading represents the current state of the sensor
type Reading struct {
	ReadableDate string `json:"ReadableDate,omitempty"`
	IsDry        bool
}

// Start read data from database
func Start(projectName string) {
	client, ctx, err := database.GetClient(projectName)
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
