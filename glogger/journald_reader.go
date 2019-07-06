package glogger

import (
    "io"
    "os/exec"
)

type JournaldReader struct {
    Reader io.Reader
    cmd *exec.Cmd
}

func NewJournaldReader() (*JournaldReader, error) {
    var err error
    reader := new(JournaldReader)
    reader.cmd = exec.Command("journalctl", "-n", "10", "--output", "json")

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
