package main

import (
	"fmt"
	"io"
	"os"

	"github.com/frickiericker/bedwinbin/bedgraph"
	"github.com/frickiericker/bedwinbin/taskch"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// Config holds run-time app configurations.
type Config struct {
	binSize int64
	input   io.Reader
	output  io.Writer
}

func run() error {
	config := Config{
		binSize: 100,
		input:   os.Stdin,
		output:  os.Stdout,
	}

	tasks := taskch.New()

	records := make(chan bedgraph.Record)
	bins := make(chan bedgraph.Record)

	tasks.Go(func() error {
		defer close(records)
		return bedgraph.Scan(config.input, records)
	})

	tasks.Go(func() error {
		defer close(bins)
		return rebin(records, bins, config.binSize)
	})

	tasks.Go(func() error {
		return bedgraph.Dump(bins, config.output)
	})

	return tasks.Wait()
}
