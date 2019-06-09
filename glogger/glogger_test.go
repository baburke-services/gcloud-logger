package glogger

import (
	"testing"
)

func TestNewJsonReaderNoSize(t *testing.T) {
	var reader *JsonReader
	reader = NewJsonReader(-1)

	if cap(reader.events) != EVENT_CHANNEL_SIZE {
		t.Errorf("event channel size: %d", cap(reader.events))
	}
}

func TestNewJsonReaderSize(t *testing.T) {
	var reader *JsonReader
	reader = NewJsonReader(10)

	if cap(reader.events) != 10 {
		t.Errorf("event channel size: %d", cap(reader.events))
	}
}

// vim: noexpandtab
