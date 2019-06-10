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
	bytes := []byte(`[
		"whatever test",
		34,
		null
	]`)
	ch := reader.FromJsonBytes(bytes)

	data, ok := <-ch
	if !ok {
		t.Errorf("unable to read from channel")
	}

	array, ok := data.([]interface{})
	if !ok {
		t.Error("returned data is not an array\n")
	}

	if len(array) != 3 {
		t.Error("array is not length 3")
	}

	s, ok := array[0].(string)
	if !ok {
		t.Error("first element is not string")
	}

	if s != "whatever test" {
		t.Errorf("got back the wrong string: %s\n", s)
	}

	_, ok = <-ch
	if ok {
		t.Error("there was more than one item on the channel")
	}
}

// vim: noexpandtab
