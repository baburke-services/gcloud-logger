package main

import (
	"testing"
)

func TestReaderNonJSON(t *testing.T) {
	command := []string{"echo", "hello"}
	ch := create_reader(command)
	_, ok := <-ch
	if ok {
		t.Fail()
	}
}

func TestReaderJsonMap(t *testing.T) {
	command := []string{"echo", "{\"output\": \"hello\"}"}
	ch := create_reader(command)
	blob := <-ch
	if _, ok := blob.(map[string]interface{}); !ok {
		t.Fail()
	}
}

func TestReaderJsonArray(t *testing.T) {
	command := []string{"echo", "[\"output\", \"hello\"]"}
	ch := create_reader(command)
	blob := <-ch
	if _, ok := blob.([]interface{}); !ok {
		t.Fail()
	}
}

func TestReaderJsonNull(t *testing.T) {
	command := []string{"echo", "null"}
	ch := create_reader(command)
	blob, ok := <-ch
	if !ok || blob != nil {
		t.Fail()
	}
}

func TestReaderJsonNumber(t *testing.T) {
	command := []string{"echo", "2"}
	ch := create_reader(command)
	blob := <-ch
	if _, ok := blob.(float64); !ok {
		t.Fail()
	}
}

func TestReaderJsonBool(t *testing.T) {
	command := []string{"echo", "false"}
	ch := create_reader(command)
	blob := <-ch
	if _, ok := blob.(bool); !ok {
		t.Fail()
	}
}
