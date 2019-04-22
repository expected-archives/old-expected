package docker

import "time"

type Output string

const (
	OuputStdout Output = "STDOUT"
	OuputStderr Output = "STDERR"
)

type LogEntry struct {
	Output  Output
	Message string
	Time    *time.Time
}

func convertEnv(env map[string]string) []string {
	var converted []string
	for k, v := range env {
		converted = append(converted, k+"="+v)
	}
	return converted
}

func ParseLogEntry(b []byte) (*LogEntry, error) {
	return nil, nil
}
