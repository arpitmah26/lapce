package lsp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os/exec"
)

// StdinoutStream is
type StdinoutStream struct {
	in     io.WriteCloser
	out    io.ReadCloser
	reader *bufio.Reader
}

// NewStdinoutStream creates
func NewStdinoutStream(command string, arg ...string) (*StdinoutStream, error) {
	cmd := exec.Command(command, arg...)
	inw, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	outr, err := cmd.StdoutPipe()
	if err != nil {
		inw.Close()
		return nil, err
	}

	err = cmd.Start()
	if err != nil {
		return nil, err
	}
	return &StdinoutStream{
		in:     inw,
		out:    outr,
		reader: bufio.NewReaderSize(outr, 8192),
	}, nil
}

// WriteObject implements ObjectStream.
func (s *StdinoutStream) WriteObject(obj interface{}) error {
	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	log.Println(string(data))
	s.in.Write([]byte(fmt.Sprintf("Content-Length: %d\r\n\r\n", len(data))))
	_, err = s.in.Write(data)
	return err
}

// ReadObject implements ObjectStream.
func (s *StdinoutStream) ReadObject(v interface{}) error {
	s.reader.ReadSlice('\n')
	s.reader.ReadSlice('\n')
	decoder := json.NewDecoder(s.reader)
	err := decoder.Decode(v)
	return err
}

// Close implements ObjectStream.
func (s *StdinoutStream) Close() error {
	err := s.in.Close()
	if err != nil {
		return err
	}
	return s.out.Close()
}