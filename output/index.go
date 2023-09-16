package rssoutput

import (
	"fmt"
	"io"
	"os"
	rssfunctions "rssreader/functions"
)

// MultiWriter represents a writer that can write to multiple underlying writers.
type MultiWriter struct {
	writers []io.Writer
}

// Write writes the provided byte slice to all the underlying writers in MultiWriter.
func (mw MultiWriter) Write(p []byte) (n int, err error) {
	for _, w := range mw.writers {
		n, err = w.Write(p)
		if err != nil {
			return n, err
		}
	}
	return n, err
}

// New creates a new MultiWriter based on the provided options.
// It can include "stdout" to write to standard output and "file" to write to a file.
func New(options []string) MultiWriter {

	var writers []io.Writer

	// Check if "stdout" is in options and add os.Stdout as a writer if found.
	if rssfunctions.Contains(options, "stdout") {
		writers = append(writers, os.Stdout)
	}

	// Check if "file" is in options and add a file writer if found.
	if rssfunctions.Contains(options, "file") {
		file, err := os.Create("rssreader.log")
		if err != nil {
			fmt.Println("Error creating file:", err)
		} else {
			writers = append(writers, file)
		}
	}

	// Create a new MultiWriter with the selected writers.
	mw := MultiWriter{
		writers: writers,
	}

	return mw
}
