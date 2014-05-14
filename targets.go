package slog

import (
	"log"
	"os"
)

const bufferSize = 100

var Stdout chan<- string
var Stderr chan<- string

func init() {
	Stdout = FileWriter(os.Stdout)
	Stderr = FileWriter(os.Stderr)
}

func FileWriter(f *os.File) chan<- string {
	ch := make(chan string, bufferSize)
	go func() {
		for line := range ch {
			_, err := f.WriteString(line)
			if err != nil {
				// Using slog to log errors about slog seems...
				// unwise.
				log.Println("slog: while writing %q to %v: %v",
					line, f, err)
			}
		}
	}()
	return ch
}

func stdout(line map[string]interface{}) {
	Stdout <- Format(line)
}
