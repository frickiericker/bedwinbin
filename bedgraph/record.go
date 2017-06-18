package bedgraph

import (
	"fmt"
)

// Record models the content of a single record in bedgraph file.
type Record struct {
	Group      string
	Start, End int64
	Value      float64
}

// IsValid determines if a record designates a valid interval.
func (rec *Record) IsValid() bool {
	return rec.Start < rec.End
}

// String formats a record as a human-readable string.
func (rec *Record) String() string {
	return fmt.Sprintf("%s %d %d %g", rec.Group, rec.Start, rec.End, rec.Value)
}
