package main

import (
	"cloud.google.com/go/logging"
	"encoding/json"
	"golang.org/x/net/context"
	"log"
	"os"
	"strconv"
)

func syslog_to_logger_level(level string) (logging.Severity, error) {
	var LEVELS = [8]logging.Severity{
		logging.Emergency,
		logging.Alert,
		logging.Critical,
		logging.Error,
		logging.Warning,
		logging.Notice,
		logging.Info,
		logging.Debug,
	}

	i, err := strconv.Atoi(level)
	if err != nil {
		return logging.Default, err
	}

	return LEVELS[i], nil
}

func main() {
	dec := json.NewDecoder(os.Stdin)
	ctx := context.Background()
	projectID := "baburke-services"

	client, err := logging.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	logName := "journald"
	lg := client.Logger(logName)

	for dec.More() {
		var dat map[string]interface{}
		dec.Decode(&dat)
		var level logging.Severity = logging.Default
		var err error

		if v, present := dat["PRIORITY"]; present {
			s, ok := v.(string)
			if ok {
				level, err = syslog_to_logger_level(s)
				if err != nil {
					log.Printf("failed to parse %s; continuing with default severity", s)
				}
			} else {
				log.Print("priority is not a string:", v)
			}
		}

		lg.Log(logging.Entry{Payload: dat, Severity: level})
	}
}

// vim: noexpandtab
