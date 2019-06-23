package glogger

import (
    "bytes"
    "testing"
)

func TestReaderJsonMap(t *testing.T) {
    json_bytes := []byte("{\"output\": \"hello\"}")
    bytes_reader := bytes.NewReader(json_bytes)
    log_reader := NewLogReader(bytes_reader)
    entry := log_reader.read_entry()
    v := entry.raw_data.(map[string]string)

	if m, present := v["output"]; !present {
        t.Fail()
    } else if m != "hello" {
        t.Fail()
    }
}

// vim: noexpandtab
