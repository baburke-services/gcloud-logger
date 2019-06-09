package glogger

const (
	EVENT_CHANNEL_SIZE = 25
)

type JsonReader struct {
	events <-chan interface{}
}

func NewJsonReader(channel_size int) *JsonReader {
	if channel_size < 0 {
		channel_size = EVENT_CHANNEL_SIZE
	}

	events := make(chan interface{}, channel_size)

	reader := JsonReader{events: events}

	return &reader
}

// vim: noexpandtab
