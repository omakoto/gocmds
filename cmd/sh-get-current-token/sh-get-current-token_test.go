package main

import (
	"fmt"
	"github.com/omakoto/go-common/src/shell"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFindToken(t *testing.T) {
	inputs := []struct {
		source       string
		pos          int
		full         bool
		unescape     bool
		expectedWord string
		expectedOk   bool
	}{
		{"abc  def", 0, false, false, "", true},
		{"abc  def", 0, true, false, "abc", true},
		{"abc  def", 1, false, false, "a", true},
		{"abc  def", 1, true, false, "abc", true},
		{"abc  def", 2, false, false, "ab", true},
		{"abc  def", 2, true, false, "abc", true},
		{"abc  def", 3, false, false, "abc", true},
		{"abc  def", 3, true, false, "abc", true},

		{"abc  def", 4, false, false, "", false},
		{"abc  def", 4, true, false, "", false},

		{"abc  def", 5, false, false, "", true},
		{"abc  def", 5, true, false, "def", true},
		{"abc  def", 6, false, false, "d", true},
		{"abc  def", 6, true, false, "def", true},
		{"abc  def", 7, false, false, "de", true},
		{"abc  def", 7, true, false, "def", true},
		{"abc  def", 8, false, false, "def", true},
		{"abc  def", 8, true, false, "def", true},

		{" abc  def", 0, false, false, "", false},
		{" abc  def", 0, true, false, "", false},

		{`"abc"\ 'def'`, 0, true, false, `"abc"\ 'def'`, true},
		{`"abc"\ 'def'`, 0, true, true, `abc def`, true},
	}
	for _, v := range inputs {
		tokens := shell.SplitToTokens(v.source)
		actual, ok := findToken(tokens, v.pos, v.full, v.unescape)

		msg := fmt.Sprintf("Source=%s Pos=%d Full=%v Unescape=%v", v.source, v.pos, v.full, v.unescape)
		assert.Equal(t, v.expectedOk, ok, msg)
		assert.Equal(t, v.expectedWord, actual, msg)
	}
}
