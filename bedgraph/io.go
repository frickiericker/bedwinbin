package bedgraph

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const numFields = 4
const delimiter = "\t"

func firstErr(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}

// Scan reads Record values from input and pushes them to specified channel.
func Scan(input io.Reader, records chan<- Record) error {
	scanner := bufio.NewScanner(input)

	for scanner.Scan() {
		record := strings.Split(scanner.Text(), delimiter)

		if n := len(record); n != numFields {
			return fmt.Errorf("expected %d tab-separated fields, but got %d", numFields, n)
		}

		group := record[0]
		start, err1 := strconv.ParseInt(record[1], 10, 64)
		end, err2 := strconv.ParseInt(record[2], 10, 64)
		value, err3 := strconv.ParseFloat(record[3], 64)

		if err := firstErr(err1, err2, err3); err != nil {
			return err
		}

		records <- Record{
			Group: group,
			Start: start,
			End:   end,
			Value: value,
		}
	}

	return scanner.Err()
}

// Dump writes Record values to output.
func Dump(records <-chan Record, output io.Writer) error {
	writer := bufio.NewWriter(output)

	for rec := range records {
		_, err := fmt.Fprintf(
			writer,
			"%s\t%d\t%d\t%.6g\n",
			rec.Group,
			rec.Start,
			rec.End,
			rec.Value,
		)
		if err != nil {
			return err
		}
	}

	return writer.Flush()
}
