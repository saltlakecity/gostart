package main

import (
	"fmt"
	"os"

	"github.com/saltlakecity/gostart/internal/commands"
	"github.com/saltlakecity/gostart/internal/styles"
)

func main() {
	if len(os.Args) < 2 {
		printUsage("")
		return
	}

	command := os.Args[1]

	switch command {
	case "create":
		if err := commands.RunCreate(); err != nil {
			fmt.Fprintf(
				os.Stderr,
				"\n%s %s\n",
				styles.Error(styles.ErrorIcon),
				styles.Error(err.Error()),
			)
			os.Exit(1)
		}
	default:
		printUsage(fmt.Sprintf("Unknown command %q", command))
	}
}

func printUsage(message string) {
	if message != "" {
		fmt.Printf(
			"\n%s %s\n",
			styles.Error(styles.ErrorIcon),
			styles.Error(message),
		)
	}

	fmt.Printf(
		"\n%s\n%s\n\n%s\n  %s\n",
		styles.Bold(styles.Accent("GoStart")),
		styles.Muted("Create a new Go project"),
		styles.Bold("Usage:"),
		styles.Accent("gostart create"),
	)
}
