package helpers

import "strings"

func SplitCommand(command string) (splitCmd []string) {
	splitCmd = strings.Split(strings.TrimSuffix(command, "\n"), " ")
	return splitCmd
}
