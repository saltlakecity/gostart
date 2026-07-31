package prompt

import (
	"bufio"
	"fmt"
	"slices"
	"strings"

	"github.com/saltlakecity/gostart/internal/styles"
)

func Ask(scanner *bufio.Scanner, question string) string {
	fmt.Printf(
		"%s %s: ",
		styles.Accent(styles.PromptIcon),
		styles.Bold(question),
	)

	scanner.Scan()

	return strings.TrimSpace(scanner.Text())
}

func Confirm(scanner *bufio.Scanner, question string) bool {
	for {
		answer := Ask(scanner, question+" "+styles.Muted("[y/N]"))
		answer = strings.ToLower(strings.TrimSpace(answer))

		switch answer {
		case "y", "yes":
			return true

		case "", "n", "no":
			return false

		default:
			fmt.Printf(
				"  %s %s\n",
				styles.Error(styles.ErrorIcon),
				styles.Error("Enter y or n"),
			)
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
		fmt.Printf("\n%s\n", styles.Bold(options.Title))

		for _, option := range options.Options {
			fmt.Printf("  %s %s\n", styles.Muted("•"), option)
		}

		question := options.Prompt
		if question == "" {
			question = options.Title
		}
		if options.Default != "" {
			question += " " + styles.Muted("["+options.Default+"]")
		}

		answer := Ask(scanner, question)
		answer = strings.ToLower(strings.TrimSpace(answer))
		if answer == "" {
			return options.Default
		}

		result := slices.Contains(options.Options, answer)
		if result {
			return answer
		}

		fmt.Printf(
			"  %s %s\n",
			styles.Error(styles.ErrorIcon),
			styles.Error(options.InvalidMessage),
		)
	}
}
