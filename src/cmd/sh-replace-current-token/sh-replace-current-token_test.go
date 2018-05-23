package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFindToken_insert(t *testing.T) {
	inputs := []struct {
		source string
		pos    int

		newWord string

		unescape bool

		expected    string
		expectedPos int
	}{
		{`abc def`, 0, `XYZ`, false, `XYZabc def`, 3},
		{`abc def`, 1, `XYZ`, false, `aXYZbc def`, 4},
		{`abc def`, 7, `XYZ`, false, `abc defXYZ`, 10},
		{`abc def`, 8, `XYZ`, false, `abc def XYZ`, 11},
		{`abc def`, 9, `XYZ`, false, `abc def XYZ`, 11},

		{`abc def`, 0, `X Z`, false, `X Zabc def`, 3},
		{`abc def`, 0, `X Z`, true, `'X Z'abc def`, 5},
	}
	for _, v := range inputs {
		result, pos := doTransform(v.source, v.pos, v.newWord, true, v.unescape)

		msg := fmt.Sprintf("Source=%s Pos=%d Word=%s Unescape=%v", v.source, v.pos, v.newWord, v.unescape)
		assert.Equal(t, v.expected, result, msg)
		assert.Equal(t, v.expectedPos, pos, msg)
	}
}

func TestFindToken_replace(t *testing.T) {
	inputs := []struct {
		source string
		pos    int

		newWord string

		unescape bool

		expected    string
		expectedPos int
	}{
		{`abc def`, 0, `XY`, false, `XY def`, 2},
		{`abc def`, 1, `XY`, false, `XY def`, 2},
		{`abc def`, 2, `XY`, false, `XY def`, 2},
		{`abc def`, 3, `XY`, false, `XY def`, 2},
		{`abc def`, 4, `XY`, false, `abc XY`, 6},
		{`abc def`, 5, `XY`, false, `abc XY`, 6},
		{`abc def`, 6, `XY`, false, `abc XY`, 6},
		{`abc def`, 7, `XY`, false, `abc XY`, 6},
		{`abc def`, 8, `XY`, false, `abc def XY`, 10},
		{`abc def`, 9, `XY`, false, `abc def XY`, 10},
		{`abc    def`, 4, `XY`, false, `abc XY   def`, 6},
		{`abc    def`, 5, `XY`, false, `abc  XY  def`, 7},
		{`abc    def`, 6, `XY`, false, `abc   XY def`, 8},
		{`abc    def`, 7, `XY`, false, `abc    XY`, 9},
		{`abc def`, 1, `X  Y`, false, `X  Y def`, 4},
		{`abc def`, 1, `X  Y`, true, `'X  Y' def`, 6},
	}
	for _, v := range inputs {
		result, pos := doTransform(v.source, v.pos, v.newWord, false, v.unescape)

		msg := fmt.Sprintf("Source=%s Pos=%d Word=%s Unescape=%v", v.source, v.pos, v.newWord, v.unescape)
		assert.Equal(t, v.expected, result, msg)
		assert.Equal(t, v.expectedPos, pos, msg)
	}
}
