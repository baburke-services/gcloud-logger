package glogger

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"strconv"
)

type LogFeeder func(<-chan interface{})

type LogReader struct {
	events  chan *LogEntry
	decoder *json.Decoder
	reader  io.Reader
}

type LogEntry struct {
	raw_data  interface{}
	Message   string
	Level     int
	LevelName string
	Cursor    string
}

type JournalctlReader struct {
	//process
	cursor string
}

const (
	EVENT_CHANNEL_SIZE = 25
	DEFAULT_LOG_LEVEL  = 6
)

var END_OF_STREAM = errors.New("reached end of stream")

var LOG_LEVELS = [8]string{
	"EMERGENCY",
	"ALERT",
	"CRITICAL",
	"ERROR",
	"WARNING",
	"NOTICE",
	"INFO",
	"DEBUG",
}

func NewLogReader(reader io.Reader) *LogReader {
	log_reader := new(LogReader)
	log_reader.decoder = json.NewDecoder(reader)
	log_reader.reader = reader

	return log_reader
}

func (r *LogReader) StartStream(channel_size int) <-chan *LogEntry {
	if channel_size < 0 {
		channel_size = EVENT_CHANNEL_SIZE

	} else if channel_size > 100 {
		log.Println("large channel size requests (>100), setting to 100")
		channel_size = 100
	}

	r.events = make(chan *LogEntry, channel_size)

	streamer := func() {
		defer close(r.events)

		for {
			if entry, err := r.read_entry(); err != nil {
				log.Print(err)
				break
			} else {
				r.events <- entry
			}
		}
	}

	go streamer()
	return r.events
}

func (r *LogReader) read_entry() (*LogEntry, error) {
	if !r.decoder.More() {
		return nil, END_OF_STREAM
	}

	var v map[string]string
	if err := r.decoder.Decode(&v); err != nil {
		return nil, err
	}

	entry := new(LogEntry)
	entry.raw_data = v

	entry.Message = v["MESSAGE"]
	entry.Cursor = v["__CURSOR"]

	if m, present := v["PRIORITY"]; !present {
		entry.Level = DEFAULT_LOG_LEVEL
		entry.LevelName = LOG_LEVELS[DEFAULT_LOG_LEVEL]

	} else if i, err := strconv.Atoi(m); err != nil {
		log.Printf("could not parse priority %s", m)
		entry.Level = DEFAULT_LOG_LEVEL
		entry.LevelName = LOG_LEVELS[DEFAULT_LOG_LEVEL]

	} else {
		entry.Level = i
		entry.LevelName = LOG_LEVELS[i]
	}

	return entry, nil
}

// vim: noexpandtab
