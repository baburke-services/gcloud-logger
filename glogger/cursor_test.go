package glogger

import (
	"bytes"
	"errors"
	"io"
	"testing"
)

type NullOpen struct{}

func (o *NullOpen) Read(buf []byte) (int, error) { return 0, io.EOF }
func (o *NullOpen) Close() error                 { return nil }

type BadRead struct{}

func (o *BadRead) Read(buf []byte) (int, error) {
	return 0, errors.New("test bad read")
}
func (o *BadRead) Close() error { return nil }

type GoodRead struct {
	buf *bytes.Buffer
}

func NewGoodRead() *GoodRead {
	reader := new(GoodRead)
	reader.buf = bytes.NewBufferString("test cursor")
	return reader
}
func (o *GoodRead) Read(buf []byte) (int, error) {
	return o.buf.Read(buf)
}
func (o *GoodRead) Close() error { return nil }

type TestOpener struct {
	reader io.ReadCloser
	err    error
}

func (o *TestOpener) Open() (io.ReadCloser, error) { return o.reader, o.err }

func TestReadCursor(t *testing.T) {
	t.Run("BadOpen", func(t *testing.T) {
		opener := &TestOpener{
			reader: &NullOpen{},
			err:    errors.New("test failed to open"),
		}
		_, err := read_cursor(opener)
		if err != ERROR_NO_CURSOR {
			t.Errorf("got %q; expected %q", err, ERROR_NO_CURSOR)
		}
	})

	t.Run("BadRead", func(t *testing.T) {
		opener := &TestOpener{
			reader: &BadRead{},
			err:    nil,
		}
		_, err := read_cursor(opener)
		if err != ERROR_NO_CURSOR {
			t.Errorf("got %q; expected %q", err, ERROR_NO_CURSOR)
		}
	})

	t.Run("GoodRead", func(t *testing.T) {
		opener := &TestOpener{
			reader: NewGoodRead(),
			err:    nil,
		}
		cursor, err := read_cursor(opener)
		if err != nil {
			t.Errorf("got %q; expected nil", err)
		} else if cursor != "test cursor" {
			t.Errorf("got %q; expect cursor %q", cursor, "test cursor")
		}
	})
}

// vim: noexpandtab
