package prompt

import (
	"bufio"
	"fmt"
	"strings"
)

var allowed = []string{
	"y", "Y", "n", "N",
}

func Ask(scanner *bufio.Scanner, question string) string {
	fmt.Print(question + ": ")

	scanner.Scan()

	return strings.TrimSpace(scanner.Text())
}

func Confirm(scanner *bufio.Scanner, question string, allowed []string) bool {
	answer := Ask(scanner, question+" [y/n]")

	answer = strings.ToLower(answer)

	return answer == "y" || answer == "yes"
}
