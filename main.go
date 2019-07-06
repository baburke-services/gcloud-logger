package main

import (
	"fmt"
	"source.baburke.net/baburke-services/gcloud-logger/glogger"
)

func main() {
	journald, err := glogger.NewJournaldReader()
	if err != nil {
		panic(err)
	}

	reader := glogger.NewLogReader(journald.Reader)
	log_channel := reader.StartStream(10)
	for entry := range log_channel {
		fmt.Printf("%s: %s\n", entry.LevelName, entry.Message)
	}

	if err := journald.Close(); err != nil {
		panic(err)
	}
}

// vim: noexpandtab
