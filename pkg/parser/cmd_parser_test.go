package parser

import (
	"testing"

	testify "github.com/stretchr/testify"
)

func TestParseCmdFromString(t *testing.T) {
	input := "execute command: docker stop xyz"
	expected := "docker stop xyz"

	actual, err := ParseCmdFromString(input)

	testify.NoError(t, err)
	testify.Equal(t, expected, actual)
}

func TestParseCmdFromStringNoMatch(t *testing.T) {
	input := "stop xyz"
	expected := ""

	actual, err := ParseCmdFromString(input)

	testify.NoError(t, err)
	testify.Equal(t, expected, actual)
}

func TestParseCmdFromStringEmptyInput(t *testing.T) {
	input := ""
	expected := ""

	actual, err := ParseCmdFromString(input)

	testify.NoError(t, err)
	testify.Equal(t, expected, actual)
}
