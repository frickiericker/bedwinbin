package bedgraph

import (
	"reflect"
	"testing"
)

func TestRecord_order(t *testing.T) {
	record := Record{"foo", 1, 2, 3.}
	expected := Record{
		Group: "foo",
		Start: 1,
		End:   2,
		Value: 3.,
	}
	if !reflect.DeepEqual(record, expected) {
		t.Error("wrong order", record)
	}
}

func TestRecord_IsValid_validCases(t *testing.T) {
	testCases := []Record{
		{Start: 0, End: 1},
	}
	for _, testCase := range testCases {
		if !testCase.IsValid() {
			t.Error("unexpectedly invalid:", testCase)
		}
	}
}

func TestRecord_IsValid_invalidCases(t *testing.T) {
	testCases := []Record{
		{Start: 0, End: 0},
		{Start: 1, End: 0},
	}
	for _, testCase := range testCases {
		if testCase.IsValid() {
			t.Error("unexpectedly valid:", testCase)
		}
	}
}

func TestRecord_String(t *testing.T) {
	testCases := []struct {
		record   Record
		expected string
	}{
		{
			record:   Record{"foo", 1, 2, 3.},
			expected: "foo 1 2 3",
		},
	}
	for _, testCase := range testCases {
		result := testCase.record.String()
		if result != testCase.expected {
			t.Error("unexpected result", result)
		}
	}
}
