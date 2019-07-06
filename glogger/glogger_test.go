package glogger

import (
	"bytes"
	"testing"
)

func TestReaderJsonMap(t *testing.T) {
	json_bytes := []byte("{\"output\": \"hello\"}")
	bytes_reader := bytes.NewReader(json_bytes)
	log_reader := NewLogReader(bytes_reader)
	entry, err := log_reader.read_entry()
	if err != nil {
		t.Fatal("failed to read a single entry")
	}

	v := entry.raw_data.(map[string]string)

	if v["output"] != "hello" {
		t.Error("raw message did not contain key-value pair")
	}

	if entry.Message != "" {
		t.Errorf("Message == %q, expected ''", entry.Message)
	}

	if entry.Cursor != "" {
		t.Errorf("Cursor == %q, expected ''", entry.Cursor)
	}

	if entry.Level != DEFAULT_LOG_LEVEL {
		t.Errorf("Log Level == %q; expected %q", entry.Level, DEFAULT_LOG_LEVEL)
	}

	if entry.LevelName != LOG_LEVELS[DEFAULT_LOG_LEVEL] {
		t.Errorf("LevelName == %q; expected %q", entry.LevelName, LOG_LEVELS[DEFAULT_LOG_LEVEL])
	}
}

// vim: noexpandtab
