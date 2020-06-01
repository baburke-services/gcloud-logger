package glogger

import (
	"errors"
	"io"
	"reflect"
	"testing"
)

func TestNewJournaldReader(t *testing.T) {
	t.Run(
		"base", testNewJournaldReaderArgs(
			"", false, []string{"journalctl", "--output", "json"},
		),
	)
	t.Run(
		"cursor", testNewJournaldReaderArgs(
			"something",
			false,
			[]string{"journalctl", "--output", "json", "--cursor", "something"},
		),
	)
	t.Run(
		"follow", testNewJournaldReaderArgs(
			"", true, []string{"journalctl", "--output", "json", "--follow"},
		),
	)
}

func testNewJournaldReaderArgs(
	cursor string, follow bool, expectedArgs []string,
) func(*testing.T) {
	testingFunction := func(t *testing.T) {
		reader := NewJournaldReader(cursor, follow)
		actualArgs := reader.command.Args()
		if !reflect.DeepEqual(actualArgs, expectedArgs) {
			t.Log("NewJournaldReader(): bad command args")
			t.Log("actual:", actualArgs)
			t.Log("expected:", expectedArgs)
			t.Fail()
		}
	}

	return testingFunction
}

type mockJRCommand struct {
	stdout_error error
	start_error  error
	wait_error   error
	kill_error   error
	args         []string
}

func (c *mockJRCommand) StdoutPipe() (io.ReadCloser, error) {
	return nil, c.stdout_error
}
func (c *mockJRCommand) Start() error {
	return c.start_error
}
func (c *mockJRCommand) Wait() error {
	return c.wait_error
}
func (c *mockJRCommand) Kill() error {
	return c.kill_error
}
func (c *mockJRCommand) Args() []string {
	return c.args
}

func TestJRStart(t *testing.T) {
	knownError := errors.New("test error")
	t.Run("base", testJRStart(&mockJRCommand{}, nil))
	t.Run(
		"stdout fail", testJRStart(
			&mockJRCommand{stdout_error: knownError}, knownError,
		),
	)
	t.Run(
		"start fail", testJRStart(
			&mockJRCommand{start_error: knownError}, knownError,
		),
	)
}

func testJRStart(command execCommand, expectedError error) func(*testing.T) {
	testingFunction := func(t *testing.T) {
		reader := &JournaldReader{
			command: command,
			reader:  nil,
		}
		actualError := reader.Start()
		if actualError != expectedError {
			t.Errorf("Start(): returned %q, expected %q", actualError, expectedError)
		}
	}
	return testingFunction
}

type mockReader struct{}

func (r *mockReader) Read(b []byte) (int, error) { return 0, nil }

func TestJRClose(t *testing.T) {
	knownError := errors.New("test error")
	t.Run(
		"nil reader",
		testJRClose(&JournaldReader{}, ErrNoProcessStarted),
	)
	t.Run(
		"kill ignored",
		testJRClose(
			&JournaldReader{
				reader:  &mockReader{},
				command: &mockJRCommand{kill_error: knownError},
			},
			nil,
		),
	)
	t.Run(
		"wait returned",
		testJRClose(
			&JournaldReader{
				reader:  &mockReader{},
				command: &mockJRCommand{wait_error: knownError},
			},
			knownError,
		),
	)
}

func testJRClose(r *JournaldReader, expectedError error) func(*testing.T) {
	testingFunction := func(t *testing.T) {
		actualError := r.Close()
		if actualError != expectedError {
			t.Errorf("Close(): expected %q, actual %q", expectedError, actualError)
		}
	}
	return testingFunction
}
