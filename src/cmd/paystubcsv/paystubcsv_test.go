package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseDallars(t *testing.T) {
	inputs := []struct {
		source   string
		expected string
	}{
		{`0`, `0`},
		{`$0`, `0`},
		{`($0)`, `0`},
		{`($1)`, `-1`},
		{`($1,234)`, `-1234`},
		{`$1`, `1`},
		{`$1,234`, `1234`},
	}
	for _, v := range inputs {
		assert.Equal(t, v.expected, parseDollars(v.source).String(), v.source)
	}
}

func TestExtractLabeledData(t *testing.T) {
	inputs := []struct {
		line  string
		label string

		expectedFound bool
		expected      string
	}{
		{"  abc ", "xyz", false, ""},
		{"  abc ", "ab", true, "c"},
		{"  abc 1  ", "ab", true, "c 1"},
		{"  abc 1  ", "abc ", true, "1"},
		{"  abc   1", "abc ", true, "1"},
	}
	for _, v := range inputs {
		found, data := extractLabeledData(v.line, v.label)

		msg := fmt.Sprintf("%#v", v)
		assert.Equal(t, v.expectedFound, found, msg)
		assert.Equal(t, v.expected, data, msg)
	}
}

func TestExtractField(t *testing.T) {
	inputs := []struct {
		line  string
		label string
		skip  int

		expectedFound bool
		expected      string
	}{
		{"  abc ", "xyz", 0, false, ""},
		{"  abc ", "ab", 0, true, "c"},
		{"  abc 1  ", "ab", 0, true, "c"},
		{"  abc 1  ", "abc ", 0, true, "1"},
		{"  abc   1", "abc ", 0, true, "1"},
		{"  abc 12   34", "abc ", 0, true, "12"},
		{"  abc 12   34", "abc ", 1, true, "34"},
	}
	for _, v := range inputs {
		found, data := extractField(v.line, v.label, v.skip)

		msg := fmt.Sprintf("%#v", v)
		assert.Equal(t, v.expectedFound, found, msg)
		assert.Equal(t, v.expected, data, msg)
	}
}
