package glogger

import (
	"encoding/json"
	"fmt"
)

type LogFeeder func(<-chan interface{})

type LogReader struct {
	events chan interface{}
}

const (
	EVENT_CHANNEL_SIZE = 25
)

func NewLogReader(channel_size int) *LogReader {
	if channel_size < 0 {
		channel_size = EVENT_CHANNEL_SIZE
	}

	events := make(chan interface{}, channel_size)

	reader := LogReader{events: events}

	return &reader
}

func (r *LogReader) FromJsonBytes(s []byte) <-chan interface{} {
	var data interface{}

	// defining inline as a special case that still works like the others
	function := func(c chan<- interface{}) {
		if err := json.Unmarshal(s, &data); err != nil {
			fmt.Errorf("unmarshal error: %q\n", err)
		}

		c <- data
		close(c)
	}

	go function(r.events)

	return r.events
}

// vim: noexpandtab
