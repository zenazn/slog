package slog

import (
	"io"
	"log"
	"os"
)

const bufferSize = 100

// Stdout is the default target of the root Logger. For performance, it is
// fairly well buffered, and a goroutine is automatically spawned to read from
// the channel and print to stdout.
var Stdout chan<- string

func init() {
	ch := make(chan string, bufferSize)
	go stdoutWriter(ch)
	Stdout = ch
}

func stdout(line map[string]interface{}) {
	Stdout <- Format(line)
}

func stdoutWriter(ch <-chan string) {
	for line := range ch {
		_, err := os.Stdout.WriteString(line)
		// Using slog to log errors about slog seems... nwise.
		if err == io.ErrShortWrite {
			log.Printf("slog: short write of %q", line)
		} else if err != nil {
			log.Printf("slog: error writing to stdout: %v", err)
		}
	}
}
