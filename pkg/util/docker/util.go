package docker

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"time"
)

type Output string

const (
	OutputStdout     Output = "STDOUT"
	OutputStderr     Output = "STDERR"
	MaxMessageLength        = 1024 * 10
)

type LogReader struct {
	reader  io.Reader
	Output  Output
	Time    time.Time
	Labels  map[string]string
	Message string
}

func NewLogReader(r io.Reader) *LogReader {
	return &LogReader{
		reader: r,
	}
}

func (s *LogReader) Read(p []byte) (int, error) {
	buf := make([]byte, 8)

	l, err := io.ReadFull(s.reader, buf)
	if err != nil {
		return 0, err
	}

	switch buf[0] {
	case 0x1:
		s.Output = OutputStdout
	case 0x2:
		s.Output = OutputStderr
	default:
		return 0, errors.New("invalid output identifier")
	}

	bufSize := binary.BigEndian.Uint32(buf[4:8])
	if bufSize > MaxMessageLength {
		return 0, errors.New("too big message")
	}
	buf = make([]byte, bufSize)

	l, err = io.ReadFull(s.reader, buf)
	if err != nil {
		return 0, err
	}

	split := bytes.SplitN(buf, []byte(" "), 3)
	t, err := time.Parse(time.RFC3339Nano, string(split[0]))
	if err != nil {
		return 0, err
	}
	s.Time = t

	s.Labels = make(map[string]string)
	labels := bytes.Split(split[1], []byte(","))
	for _, label := range labels {
		kv := bytes.Split(label, []byte("="))
		s.Labels[string(kv[0])] = string(kv[1])
	}

	s.Message = string(split[2])
	copy(p, buf)
	return l, nil
}

func convertEnv(env map[string]string) []string {
	var converted []string
	for k, v := range env {
		converted = append(converted, k+"="+v)
	}
	return converted
}
