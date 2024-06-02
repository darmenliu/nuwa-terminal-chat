package parser

import (
	"fmt"
	"log/slog"
	"os"
	"regexp"
)

// The input string is like:
// execute command: docker stop xyz, this func just parse the command from the string
func ParseCmdFromString(input string) (string, error) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	re := regexp.MustCompile(`execute command: (.*)`)
	match := re.FindStringSubmatch(input)
	if match != nil {
		logger.Info("Matched:", "match content", match[1])
		return match[1], nil
	} else {
		logger.Info("No match found.")
		return "", fmt.Errorf("no match found")
	}
}
