package glogger

func CursorProcessor(events <-chan *LogEntry) (chan<- *LogEntry, error) {
	out_events := make(chan *LogEntry, EVENT_CHANNEL_SIZE)
	var cursor string

	// test access to cursor file

	go func() {
		defer close(out_events)
		for event := range events {
			cursor = event.Cursor
			out_events <- event
		}
	}()

	return out_events, nil
}

// vim: noexpandtab
