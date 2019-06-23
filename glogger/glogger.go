package glogger

import (
	"encoding/json"
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
		for {
			entry := r.read_entry()
			if entry == nil {
				close(r.events)
				break
			}
			r.events <- entry
		}
	}

	go streamer()
	return r.events
}

func (r *LogReader) read_entry() *LogEntry {
	if !r.decoder.More() {
		return nil
	}

	var v map[string]string
	err := r.decoder.Decode(&v)
	if err != nil {
		log.Println(err)
		return nil
	}

	entry := new(LogEntry)
	entry.raw_data = v

	if m, present := v["MESSAGE"]; present {
		entry.Message = m
	}

	if m, present := v["PRIORITY"]; !present {
		entry.Level = DEFAULT_LOG_LEVEL
		entry.LevelName = "DEFAULT"

	} else if i, err := strconv.Atoi(m); err != nil {
		log.Printf("could not parse priority %s", m)
		entry.Level = DEFAULT_LOG_LEVEL
		entry.LevelName = "DEFAULT"

	} else {
		entry.Level = i
		entry.LevelName = LOG_LEVELS[i]
	}

	if m, present := v["__CURSOR"]; present {
		entry.Cursor = m
	}

	return entry
}

// vim: noexpandtab
