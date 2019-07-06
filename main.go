package main

import (
	"fmt"
	"os"
	"source.baburke.net/baburke-services/gcloud-logger/glogger"
)

func main() {
	reader := glogger.NewLogReader(os.Stdin)
	log_channel := reader.StartStream(10)
	for entry := range log_channel {
		if entry == nil {
			break
		}
		fmt.Printf("%s: %s\n", entry.LevelName, entry.Message)
	}

	return
}

// vim: noexpandtab
