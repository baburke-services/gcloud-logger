package glogger

import (
	"log"
	"os"
)

func new_logger(candidate *log.Logger) *log.Logger {
	if candidate != nil {
		return candidate
	}

	logger := log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lshortfile)

	return logger
}
