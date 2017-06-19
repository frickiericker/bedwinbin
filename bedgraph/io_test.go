package bedgraph

import (
	"bytes"
	"reflect"
	"strings"
	"testing"

	"github.com/frickiericker/acol/iomock"
)

func pop(records <-chan Record) []Record {
	all := []Record{}
	for {
		select {
		case rec := <-records:
			all = append(all, rec)
		default:
			return all
		}
	}
}

func TestScan_parsesValidInputs(t *testing.T) {
	testCases := []struct {
		input    string
		expected []Record
	}{
		{
			input:    "",
			expected: []Record{},
		},
		{
			input: "foo\t0\t10\t42\n",
			expected: []Record{
				Record{"foo", 0, 10, 42.0},
			},
		},
		{
			input: "foo\t0\t10\t42\n" + "foo\t10\t20\t100\n",
			expected: []Record{
				{"foo", 0, 10, 42.},
				{"foo", 10, 20, 100.},
			},
		},
		{
			input: "foo\t10\t20\t42\n" + "bar\t0\t10\t100\n",
			expected: []Record{
				{"foo", 10, 20, 42.0},
				{"bar", 0, 10, 100.0},
			},
		},
	}

	for _, testCase := range testCases {
		input := strings.NewReader(testCase.input)
		records := make(chan Record, len(testCase.expected))

		if err := Scan(input, records); err != nil {
			t.Errorf("unexpected error < %s > parsing: %s", err, testCase.input)
		}

		result := pop(records)
		if !reflect.DeepEqual(result, testCase.expected) {
			t.Errorf("unexpected result < %v > parsing: %s", result, testCase.input)
		}
	}
}

func TestScan_failsOnInvalidInput(t *testing.T) {
	testCases := []string{
		"\n",
		"foo\n",
		"foo\t0\n",
		"foo\t0\t10\n",
		"foo\t0\t10\t\n",
		"foo\tbar\t10\t1.0\n",
		"foo\t0\tbaz\t1.0\n",
		"foo\t0\t10\tqux\n",
	}

	for _, testCase := range testCases {
		input := strings.NewReader(testCase)
		records := make(chan Record, 1)

		if err := Scan(input, records); err == nil {
			t.Error("unexpected success on parsing: %s", testCase)
		}
	}
}

func TestScan_propagateIOError(t *testing.T) {
	input := new(iomock.FailingIO)
	records := make(chan Record, 1)
	if err := Scan(input, records); err != iomock.Error {
		t.Error("unexpected error:", err)
	}
}

func TestDump(t *testing.T) {
	testCases := []struct {
		input    []Record
		expected string
	}{
		{
			input:    []Record{},
			expected: "",
		},
		{
			input: []Record{
				{"foo", 0, 10, 1.},
			},
			expected: "foo\t0\t10\t1\n",
		},
		{
			input: []Record{
				{"foo", 0, 10, 1.},
				{"foo", 10, 20, 2.},
			},
			expected: "foo\t0\t10\t1\n" + "foo\t10\t20\t2\n",
		},
		{
			input: []Record{
				{"foo", 0, 10, 1.},
				{"bar", 1, 11, 2.},
				{"baz", 2, 12, 3.},
			},
			expected: "foo\t0\t10\t1\n" + "bar\t1\t11\t2\n" + "baz\t2\t12\t3\n",
		},
	}

	for _, testCase := range testCases {
		records := make(chan Record, len(testCase.input))
		for _, rec := range testCase.input {
			records <- rec
		}
		close(records)

		output := new(bytes.Buffer)
		if err := Dump(records, output); err != nil {
			t.Error("unexpected error:", err)
		}

		result := output.String()
		if result != testCase.expected {
			t.Error("unexpected result:", result)
		}
	}
}

func TestDump_propagateIOError(t *testing.T) {
	records := make(chan Record, 1)
	records <- Record{}
	close(records)

	output := new(iomock.FailingIO)
	if err := Dump(records, output); err != iomock.Error {
		t.Error("unexpected error:", err)
	}
}
