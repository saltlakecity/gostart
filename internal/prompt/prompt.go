package prompt

import (
	"bufio"
	"fmt"
	"slices"
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

type ChoiceOptions struct {
	Title          string
	Prompt         string
	Options        []string
	Default        string
	InvalidMessage string
}

func Choice(scanner *bufio.Scanner, options *ChoiceOptions) string {
	for {
		for _, option := range options.Options {
			fmt.Println(option)
		}

		answer := Ask(scanner, options.Title)
		answer = strings.ToLower(strings.TrimSpace(answer))
		if answer == "" {
			return options.Default
		}

		result := slices.Contains(options.Options, answer)
		if result {
			return answer
		}

		fmt.Println(options.InvalidMessage)
	}

}
