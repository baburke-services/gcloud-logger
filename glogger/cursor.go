package glogger

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
)

const CURSOR_STATE_FILE = "/var/run/glogger.journald.cursor"

var ERROR_NO_CURSOR = errors.New("cannot read cursor")

func ReadCurrentCursor() (string, error) {
	var cursor bytes.Buffer

	file, err := os.Open(CURSOR_STATE_FILE)
	if err != nil {
		log.Printf("%v: cannot open cursor file", err)
		return "", ERROR_NO_CURSOR
	}
	defer file.Close()

	read, err := cursor.ReadFrom(file)
	if err != nil {
		log.Printf("%v: cannot read cursor file", err)
		return "", ERROR_NO_CURSOR
	}
	log.Printf("read %v from %v", read, CURSOR_STATE_FILE)

	return cursor.String(), nil
}

func CursorProcessor(events <-chan *LogEntry) (<-chan *LogEntry, error) {
	var file *os.File

	out_events := make(chan *LogEntry, EVENT_CHANNEL_SIZE)

	file, err := os.OpenFile(CURSOR_STATE_FILE, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}

	go func() {
		defer file.Close()
		defer close(out_events)

		for event := range events {
			cursor := bytes.NewBufferString(event.Cursor)
			if err := write_cursor(file, cursor.Bytes()); err != nil {
				log.Print(err)
				break
			}
			out_events <- event
		}
	}()

	return out_events, nil
}

// important that the file referred to here be on fast media, e.g. in memory,
// as it flushes on every write. otherwise, it should be timed to only write
// infrequently
func write_cursor(file *os.File, cursor []byte) error {
	var written int

	if _, err := file.Seek(0, 0); err != nil {
		return err
	}

	if err := file.Truncate(0); err != nil {
		return err
	}

	written, err := file.Write(cursor)
	if err != nil {
		return err
	}

	len_cursor := len(cursor)
	if written != len_cursor {
		return fmt.Errorf("wrote %v bytes; expected %v", written, len_cursor)
	}

	if err := file.Sync(); err != nil {
		return err
	}

	return nil
}

// vim: noexpandtab
