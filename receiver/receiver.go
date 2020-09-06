package receiver

import (
	"log"
	"sort"

	"github.com/kevinnayar/aggregarian/database"
)

// Reading represents the current state of the sensor
type Reading struct {
	ReadableDate string `json:"ReadableDate,omitempty"`
	IsDry        bool
}

// ReadingResult formattted with key for arrays/slices
type ReadingResult struct {
	UTCISODate   string `json:"UTCISODate,omitempty"`
	ReadableDate string `json:"ReadableDate,omitempty"`
	IsDry        bool
}

func getData(projectName string) ([]ReadingResult, error) {
	client, ctx, err := database.GetClient(projectName)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	reading := map[string]Reading{}
	if err := client.NewRef("log/").Get(ctx, &reading); err != nil {
		log.Fatalln("Error getting value:", err)
		return nil, err
	}

	keys := make([]string, 0, len(reading))
	for k := range reading {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	result := []ReadingResult{}
	for _, k := range keys {
		value := reading[k]
		result = append(result, ReadingResult{
			UTCISODate:   k,
			ReadableDate: value.ReadableDate,
			IsDry:        value.IsDry,
		})
	}

	return result, nil
}

// GetAll read all data from database
func GetAll(projectName string) ([]ReadingResult, error) {
	data, err := getData(projectName)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// GetLatest read latest entry from database
func GetLatest(projectName string) (ReadingResult, error) {
	latest := ReadingResult{}
	data, err := getData(projectName)
	if err != nil {
		return data[0], err
	}

	latest = data[len(data)-1]

	return latest, nil
}
