package main

import (
	"fmt"

	"github.com/frickiericker/bedwinbin/bedgraph"
)

func rebin(records <-chan bedgraph.Record, bins chan<- bedgraph.Record, binSize int64) error {

	normalize := func(bin bedgraph.Record) bedgraph.Record {
		bin.Value /= float64(bin.End - bin.Start)
		return bin
	}

	bin := bedgraph.Record{}

	for rec := range records {
		if bin.Group != rec.Group {
			if bin.IsValid() {
				bins <- normalize(bin)
			}

			bin = bedgraph.Record{
				Group: rec.Group,
				Start: 0,
				End:   binSize,
				Value: 0,
			}
		}

		if rec.Start < bin.Start {
			return fmt.Errorf("out of order input: %s", rec.String())
		}

		for bin.End <= rec.Start {
			bins <- normalize(bin)

			bin.Start += binSize
			bin.End += binSize
			bin.Value = 0
		}

		for bin.End < rec.End {
			span := float64(bin.End - rec.Start)
			bin.Value += span * rec.Value
			bins <- normalize(bin)

			rec.Start = bin.End

			bin.Start += binSize
			bin.End += binSize
			bin.Value = 0
		}

		span := float64(rec.End - rec.Start)
		bin.Value += span * rec.Value
	}

	if bin.IsValid() {
		bins <- normalize(bin)
	}

	return nil
}
