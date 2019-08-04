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
	args := []string{"--output", "json"}

	if cursor != "" {
		args = append(args, "--cursor", cursor)
	}

	if follow {
		args = append(args, "--follow")
	}

	reader.cmd = exec.Command(name, args...)

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
