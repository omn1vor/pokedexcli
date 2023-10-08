package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "test String",
			expected: []string{"test", "string"},
		},
		{
			input:    "  _Test_  SECOND",
			expected: []string{"_test_", "second"},
		},
	}

	for _, cs := range cases {
		actual := cleanInput(cs.input)
		if len(actual) != len(cs.expected) {
			t.Errorf("Expected the same number of words: %v, got: %v", len(cs.expected), len(actual))
		}
		for i, v := range cs.expected {
			if actual[i] != v {
				t.Errorf("Expected word #%v to be: %v, got: %v)", i, actual[i], v)
			}
		}
	}
}
