package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCmdFromString(t *testing.T) {
	input := "execute command: docker stop xyz"
	expected := "docker stop xyz"

	actual, err := ParseCmdFromString(input)

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestParseCmdFromStringNoMatch(t *testing.T) {
	input := "stop xyz"
	expected := ""

	actual, err := ParseCmdFromString(input)

	assert.EqualError(t, err, "no match found")
	assert.Equal(t, expected, actual)
}

func TestParseCmdFromStringEmptyInput(t *testing.T) {
	input := ""
	expected := ""

	actual, err := ParseCmdFromString(input)

	assert.EqualError(t, err, "no match found")
	assert.Equal(t, expected, actual)
}
