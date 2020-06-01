package glogger

import (
	"errors"
	"io"
	"os/exec"
)

var (
	ErrNoProcessStarted = errors.New("No process started")
)

type execCommand interface {
	StdoutPipe() (io.ReadCloser, error)
	Start() error
	Wait() error
	Kill() error
	Args() []string
}

type builtinCommand struct {
	command *exec.Cmd
}

func (c *builtinCommand) StdoutPipe() (io.ReadCloser, error) {
	return c.command.StdoutPipe()
}

func (c *builtinCommand) Start() error {
	return c.command.Start()
}

func (c *builtinCommand) Wait() error {
	return c.command.Wait()
}

func (c *builtinCommand) Kill() error {
	return c.command.Process.Kill()
}

func (c *builtinCommand) Args() []string {
	return c.command.Args
}

type JournaldReader struct {
	reader  io.Reader
	command execCommand
}

func NewJournaldReader(cursor string, follow bool) *JournaldReader {
	name := "journalctl"
	args := []string{"--output", "json"}

	if cursor != "" {
		args = append(args, "--cursor", cursor)
	}

	if follow {
		args = append(args, "--follow")
	}

	command := &builtinCommand{
		command: exec.Command(name, args...),
	}

	reader := &JournaldReader{
		reader:  nil,
		command: command,
	}

	return reader
}

func (r *JournaldReader) Start() error {
	reader, err := r.command.StdoutPipe()
	if err != nil {
		return err
	}
	r.reader = reader

	if err := r.command.Start(); err != nil {
		return err
	}

	return nil
}

// Close Kill any existing sub-process, wait for it, and return any error
func (r *JournaldReader) Close() error {
	if r.reader == nil {
		return ErrNoProcessStarted
	}
	r.command.Kill()
	return r.command.Wait()
}

// vim: noexpandtab
