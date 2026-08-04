package main

import (
	"fmt"
	"io"
	"os"

	"github.com/saltlakecity/gostart/internal/commands"
	"github.com/saltlakecity/gostart/internal/styles"
)

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	if len(args) == 0 {
		printUsage(os.Stdout, "")
		return 0
	}

	switch args[0] {
	case "create":
		if len(args) != 1 {
			printUsage(
				os.Stderr,
				`Command "create" does not accept arguments yet`,
			)
			return 2
		}

		if err := commands.RunCreate(); err != nil {
			fmt.Fprintf(
				os.Stderr,
				"\n%s %s\n",
				styles.Error(styles.ErrorIcon),
				styles.Error(err.Error()),
			)
			return 1
		}

		return 0

	case "help", "--help", "-h":
		printUsage(os.Stdout, "")
		return 0

	default:
		printUsage(
			os.Stderr,
			fmt.Sprintf("Unknown command %q", args[0]),
		)
		return 2
	}
}

func printUsage(writer io.Writer, message string) {
	if message != "" {
		fmt.Fprintf(
			writer,
			"\n%s %s\n",
			styles.Error(styles.ErrorIcon),
			styles.Error(message),
		)
	}

	fmt.Fprintf(
		writer,
		"\n%s\n%s\n\n%s\n  %s\n",
		styles.Bold(styles.Accent("GoStart")),
		styles.Muted("Create a new Go project"),
		styles.Bold("Usage:"),
		styles.Accent("gostart create"),
	)
}
