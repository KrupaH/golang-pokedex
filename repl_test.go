package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{{
		input:    "   hello world  ",
		expected: []string{"hello", "world"},
	}, {
		input:    " test   one ",
		expected: []string{"test", "one"},
	}}

	for _, c := range cases {
		actual := cleanInput(c.input)
		for idx, word := range actual {
			if word != c.expected[idx] {
				t.Errorf("Expected %s, got %s at index %d", c.expected[idx], word, idx)
			}
		}
	}

}
