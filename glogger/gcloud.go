package glogger

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/logging"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/compute/v1"
)

const (
	PROJECT_DEFAULT = "baburke-services"
)

func get_project_id() string {
	project := os.Getenv("GLOGGER_PROJECT")
	if project != "" {
		return project
	}

	ctx := context.Background()
	credentials, err := google.FindDefaultCredentials(ctx, compute.ComputeScope)
	if err != nil {
		log.Printf("warning: %v", err)
	}

	if credentials.ProjectID != "" {
		return credentials.ProjectID
	}

	log.Printf("warning: returning builtin default for project")
	return PROJECT_DEFAULT
}

func NewGLogger() (*logging.Logger, *logging.Client, error) {
	ctx := context.Background()
	projectID := get_project_id()

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
				Payload:  entry.raw_data,
			})
			count += 1
		}
		log.Printf("processed %d entries", count)
	}()

	return done, nil
}

// vim: noexpandtab
