package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{
			input:    "test",
			expected: "test",
		},
		{
			input:    "  _Test_  ",
			expected: "_test_",
		},
	}

	for _, cs := range cases {
		actual := cleanInput(cs.input)
		if actual != cs.expected {
			t.Errorf("Expected: %v, got: %v", cs.expected, actual)
		}
	}
}
