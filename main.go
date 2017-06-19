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
	binSize, winSize int64
	input            io.Reader
	output           io.Writer
}

func run() error {
	config := Config{
		binSize: 100,
		winSize: 90,
		input:   os.Stdin,
		output:  os.Stdout,
	}

	offset := config.binSize/2 - config.winSize/2
	unitSize := gcd(config.binSize, config.winSize)
	fmt.Fprintf(os.Stderr, "Unit:\t%d\n", unitSize)
	fmt.Fprintf(os.Stderr, "Offset:\t%d\n", offset)
	fmt.Fprintf(os.Stderr, "Bin:\t%d\n", config.binSize/unitSize)
	fmt.Fprintf(os.Stderr, "Win:\t%d\n", config.winSize/unitSize)

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

func gcd(m, n int64) int64 {
	for n > 0 {
		n, m = m%n, n
	}
	return m
}
