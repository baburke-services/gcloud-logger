package main

import (
	"fmt"
	"log"
	"source.baburke.net/baburke-services/gcloud-logger/glogger"
)

func main() {
	cursor, err := glogger.ReadCurrentCursor()
	if err == glogger.ERROR_NO_CURSOR {
		log.Print("cursor not found; proceeding")
	} else if err != nil {
		panic(err)
	}

	journald, err := glogger.NewJournaldReader(cursor)
	if err != nil {
		panic(err)
	}

	log.Printf("read cursor %q", cursor)

	reader := glogger.NewLogReader(journald.Reader)
	log_channel := reader.StartStream(0)
	post_cursor, err := glogger.CursorProcessor(log_channel)
	if err != nil {
		panic(err)
	}

	for entry := range post_cursor {
		fmt.Printf("%s: %s\n", entry.LevelName, entry.Message)
	}

	if err := journald.Close(); err != nil {
		panic(err)
	}
}

// vim: noexpandtab
