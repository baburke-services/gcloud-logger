package glogger

import (
	"errors"
	"io"
	"os/exec"
	"reflect"
	"testing"
)

func TestNewJournaldReaderNoCursorNoFollow(t *testing.T) {
	reader := NewJournaldReader("", false)
	cmd := reader.command.(*exec.Cmd)
	expected_args := []string{"journalctl", "--output", "json"}

	if !reflect.DeepEqual(cmd.Args, expected_args) {
		t.Log("bad command args")
		t.Log("actual:", cmd.Args)
		t.Log("expected:", expected_args)
		t.Fail()
	}
}

func TestNewJournaldReaderCursor(t *testing.T) {
	reader := NewJournaldReader("test_cursor", false)
	cmd := reader.command.(*exec.Cmd)
	for i, v := range cmd.Args {
		if v != "--cursor" {
			continue
		}
		if i+1 < len(cmd.Args) && cmd.Args[i+1] == "test_cursor" {
			return
		}
	}
	t.Log("cursor arguments not in args")
	t.Log("actual:", cmd.Args)
	t.FailNow()
}

func TestNewJournaldReaderFollow(t *testing.T) {
	reader := NewJournaldReader("", true)
	cmd := reader.command.(*exec.Cmd)
	for _, v := range cmd.Args {
		if v == "--follow" {
			return
		}
	}
	t.Log("follow arguments not in args")
	t.Log("actual:", cmd.Args)
	t.FailNow()
}

type mock_jr_command struct {
	stdout_error error
	start_error  error
	wait_error   error
}

func (c *mock_jr_command) StdoutPipe() (io.ReadCloser, error) {
	return nil, c.stdout_error
}
func (c *mock_jr_command) Start() error {
	return c.start_error
}
func (c *mock_jr_command) Wait() error {
	return c.wait_error
}
func TestJRStart(t *testing.T) {
	reader := new(JournaldReader)
	command := new(mock_jr_command)
	command.stdout_error = nil
	command.start_error = nil
	reader.command = command
	result := reader.Start()
	if result != nil {
		t.Error("Start() result not nil:", result)
	}
}

func TestJRStartStdoutFail(t *testing.T) {
	reader := new(JournaldReader)
	command := new(mock_jr_command)
	command.stdout_error = errors.New("mock error")
	command.start_error = nil
	reader.command = command
	result := reader.Start()
	if result != OPEN_STDOUT_FAILED {
		t.Error("Start() gave the wrong error:", result)
	}
}

func TestJRStartStartFail(t *testing.T) {
	reader := new(JournaldReader)
	command := new(mock_jr_command)
	command.stdout_error = nil
	command.start_error = errors.New("mock error")
	reader.command = command
	result := reader.Start()
	if result != EXECUTE_FAILED {
		t.Error("Start() gave the wrong error:", result)
	}
}

func TestJRCloseReturnsWait(t *testing.T) {
	return_value := errors.New("mock error")
	command := mock_jr_command{
		start_error:  nil,
		stdout_error: nil,
		wait_error:   return_value,
	}
	reader := JournaldReader{
		command: &command,
	}
	if result := reader.Close(); result != return_value {
		t.Error("Wait() gave wrong return_value:", result)
	}
}
