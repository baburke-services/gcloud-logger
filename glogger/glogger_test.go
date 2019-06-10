package glogger

import (
	"testing"
)

func TestNewLogReaderNoSize(t *testing.T) {
	var reader *LogReader
	reader = NewLogReader(-1)

	if cap(reader.events) != EVENT_CHANNEL_SIZE {
		t.Errorf("event channel size: %d", cap(reader.events))
	}
}

func TestNewLogReaderSize(t *testing.T) {
	var reader *LogReader
	reader = NewLogReader(10)

	if cap(reader.events) != 10 {
		t.Errorf("event channel size: %d", cap(reader.events))
	}
}

func TestFromJsonString(t *testing.T) {
	reader := NewLogReader(-1)
	bytes := []byte(`"whatever test"`)
	ch := reader.FromJsonBytes(bytes)
	var data interface{}
	var ok bool
	var s string

	data = <-ch
	s, ok = data.(string)
	if !ok {
		t.Errorf("returned data is not a string!\n")
	} else if s != "whatever test" {
		t.Errorf("got back the wrong string: %s\n", s)
	}
	data, ok = <-ch
	if ok {
		t.Errorf("this was not the only item on the channel!")
	}
}

// vim: noexpandtab
