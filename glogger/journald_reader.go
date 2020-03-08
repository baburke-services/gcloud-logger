package glogger

import (
	"errors"
	"io"
	"os/exec"
)

var (
	OPEN_STDOUT_FAILED = errors.New("Failed to open Stdout pipe")
	EXECUTE_FAILED     = errors.New("Failed to start command")
)

type execCommand interface {
	StdoutPipe() (io.ReadCloser, error)
	Start() error
	Wait() error
	GetArgs() []string
}

type JournaldReader struct {
	Reader  io.Reader
	command execCommand
}

type journaldExec struct {
	cmd *exec.Cmd
}

func (r *journaldExec) StdoutPipe() (io.ReadCloser, error) {
	return r.cmd.StdoutPipe()
}

func (r *journaldExec) Start() error {
	return r.cmd.Start()
}

func (r *journaldExec) Wait() error {
	return r.cmd.Wait()
}

func (r *journaldExec) GetArgs() []string {
	return r.cmd.Args
}

func NewJournaldReader(cursor string, follow bool) *JournaldReader {
	reader := new(JournaldReader)
	name := "journalctl"
	args := []string{"--output", "json"}

	if cursor != "" {
		args = append(args, "--cursor", cursor)
	}

	if follow {
		args = append(args, "--follow")
	}

	command := new(journaldExec)
	command.cmd = exec.Command(name, args...)
	reader.command = command

	return reader
}

func (r *JournaldReader) Start() error {
	logger := new_logger(nil)
	reader, err := r.command.StdoutPipe()
	if err != nil {
		logger.Print("could not open stdout pipe:", err)
		return OPEN_STDOUT_FAILED
	}
	r.Reader = reader

	if err := r.command.Start(); err != nil {
		logger.Print("failed to execute command:", err)
		return EXECUTE_FAILED
	}

	return nil
}

func (r *JournaldReader) Close() error {
	err := r.command.Wait()
	return err
}

// vim: noexpandtab
