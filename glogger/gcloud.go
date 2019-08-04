package glogger

import (
    "context"
    "cloud.google.com/go/logging"
	"log"
	"os"
)

func NewGLogger() (*logging.Logger, error) {
	ctx := context.Background()
	projectID := "baburke-services"

	client, err := logging.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	//defer client.Close()

	logName, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	log.Printf("using log name %q", logName)

	logger := client.Logger(logName)

	return logger, nil
}

func StartLogForwarder(
	logger *logging.Logger, entries <-chan *LogEntry,
) (<-chan bool, error) {
	done := make(chan bool)
	go func() {
		defer logger.Flush()
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
		close(done)
	}()

	return done, nil
}

// vim: noexpandtab
