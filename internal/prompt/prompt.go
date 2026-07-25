package prompt

import (
	"bufio"
	"fmt"
	"strings"
)

func Ask(scanner *bufio.Scanner, question string) string {
	fmt.Print(question + ": ")

	scanner.Scan()

	return strings.TrimSpace(scanner.Text())
}

func Confirm(scanner *bufio.Scanner, question string) bool {
	for {
		answer := Ask(scanner, question+" [y/n]")
		answer = strings.ToLower(strings.TrimSpace(answer))

		switch answer {
		case "y", "yes":
			return true

		case "n", "no":
			return false

		default:
			fmt.Println("Enter [y] or [n]")
		}
	}
}
