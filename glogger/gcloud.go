package glogger

import (
    "context"
    "cloud.google.com/go/logging"
	"log"
	"os"
)

func NewGLogger() (*logging.Logger, *logging.Client, error) {
	ctx := context.Background()
	projectID := "baburke-services"

	client, err := logging.NewClient(ctx, projectID)
	if err != nil {
		return nil, nil, err
	}

	logName, err := os.Hostname()
	if err != nil {
		client.Close()
		return nil, nil, err
	}

	log.Printf("using log name %q", logName)

	logger := client.Logger(logName)

	return logger, client, nil
}

func StartLogForwarder(
	logger *logging.Logger, entries <-chan *LogEntry,
) (<-chan bool, error) {
	done := make(chan bool)
	go func() {
		defer close(done)

		var count uint64 = 0
		for entry := range entries {
			severity := logging.ParseSeverity(entry.LevelName)
			logger.Log(logging.Entry{
				Severity: severity,
				Payload: entry.raw_data,
			})
			count += 1
		}
		log.Printf("processed %d entries", count)
	}()

	return done, nil
}

// vim: noexpandtab
