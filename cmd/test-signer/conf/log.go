package conf

import (
	"io"
	"log"
	"os"
	"strings"
)

func BuildLogger(target string, prefix string, flag int) *log.Logger {
	writer := buildWriter(target)
	return log.New(writer, prefix, flag)
}

func buildWriter(target string) io.Writer {
	if strings.HasPrefix(target, "file:") {
		file, err := os.OpenFile(strings.TrimPrefix(target, "file:"),
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			panic(err)
		}
		return file
	}

	switch strings.ToLower(target) {
	case "stdout":
		return os.Stdout
	case "stderr":
		return os.Stderr
	default:
		return os.Stdout
	}
}
