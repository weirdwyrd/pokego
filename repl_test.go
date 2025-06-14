package main

import "testing"

// TestCleanInput tests the input cleaning functionality
func TestCleanInput(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "hello    world",
			expected: []string{"hello", "world"},
		},
		{
			input:    "hello world  ",
			expected: []string{"hello", "world"},
		},
	}

	for _, test := range tests {
		actual := cleanInput(test.input)

		for i, value := range actual {
			if value != test.expected[i] {
				t.Errorf("Expected %s, but got %s", test.expected[i], value)
			}
		}
	}
}
