package glogger

import (
	"io"
	"os/exec"
)

type JournaldReader struct {
	Reader io.Reader
	cmd    *exec.Cmd
}

func NewJournaldReader(cursor string, follow bool) (*JournaldReader, error) {
	var err error
	reader := new(JournaldReader)
	name := "journalctl"
	command := []string{"journalctl", "--output", "json"}

	if cursor != "" {
		command = append(command, "--cursor", cursor)
	}

	if follow {
		command = append(command, "--follow")
	}

	reader.cmd = exec.Command(name, command...)

	reader.Reader, err = reader.cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := reader.cmd.Start(); err != nil {
		return nil, err
	}

	return reader, nil
}

func (r *JournaldReader) Close() error {
	err := r.cmd.Wait()
	return err
}

// vim: noexpandtab
