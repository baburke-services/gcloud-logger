package glogger

import (
	"bytes"
	"testing"
)

func TestReadEntryDefaults(t *testing.T) {
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

func TestReadEntryCursor(t *testing.T) {
	json_bytes := []byte(`{
		"__CURSOR": "test cursor point"
	}`)
	bytes_reader := bytes.NewReader(json_bytes)
	log_reader := NewLogReader(bytes_reader)
	entry, err := log_reader.read_entry()
	if err != nil {
		t.Fatal("failed to read a single entry")
	}

	if entry.Cursor != "test cursor point" {
		t.Errorf("Cursor == %q; expected %q", entry.Cursor, "test cursor point")
	}
}

func TestReadEntryMessage(t *testing.T) {
	json_bytes := []byte(`{
		"MESSAGE": "test message"
	}`)
	bytes_reader := bytes.NewReader(json_bytes)
	log_reader := NewLogReader(bytes_reader)
	entry, err := log_reader.read_entry()
	if err != nil {
		t.Fatal("failed to read a single entry")
	}

	if entry.Message != "test message" {
		t.Errorf("Message == %q; expected %q", entry.Message, "test message")
	}
}

func TestReadEntryLogLevel(t *testing.T) {
	json_bytes := []byte(`{
		"PRIORITY": "2"
	}`)
	bytes_reader := bytes.NewReader(json_bytes)
	log_reader := NewLogReader(bytes_reader)
	entry, err := log_reader.read_entry()
	if err != nil {
		t.Fatal("failed to read a single entry")
	}

	if entry.Level != 2 {
		t.Errorf("Log Level == %q; expected %q", entry.Level, 2)
	}

	if entry.LevelName != LOG_LEVELS[2] {
		t.Errorf("Log Level Name == %q; expected %q", entry.LevelName, LOG_LEVELS[2])
	}
}

func TestStartStream(t *testing.T) {
	json_bytes := []byte(`{
		"PRIORITY": "2"
	}`)
	bytes_reader := bytes.NewReader(json_bytes)
	log_reader := NewLogReader(bytes_reader)

	ch := log_reader.StartStream(0)
	entry, ok := <-ch

	if !ok {
		t.Fatalf("no events on channel")
	}

	if entry.Level != 2 {
		t.Errorf("Log Level == %q; expected %q", entry.Level, 2)
	}

	entry, ok = <-ch
	if ok {
		t.Errorf("channel not closed")
	}

}

// vim: noexpandtab
