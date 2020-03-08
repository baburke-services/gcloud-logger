package glogger

import (
	"log"
	"os"
	"testing"
)

func TestNewLoggerGiven(t *testing.T) {
	logger := log.New(os.Stderr, "", log.LstdFlags)
	result := new_logger(logger)
	if logger != result {
		t.Fatal("returned logger was not the same as the one passed in")
	}
}

func TestNewLoggerNil(t *testing.T) {
	result := new_logger(nil)
	if result.Writer() != os.Stderr {
		t.Fatal("log writer is not stderr")
	}
	if result.Flags()&log.Lshortfile == 0 {
		t.Fatal("filename is not in the log lines")
	}
}
