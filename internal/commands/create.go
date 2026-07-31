package commands

import (
	"bufio"
	"fmt"
	"os"

	"github.com/saltlakecity/gostart/internal/generate"
	"github.com/saltlakecity/gostart/internal/models"
	"github.com/saltlakecity/gostart/internal/prompt"
	"github.com/saltlakecity/gostart/internal/styles"
)

func collectProjectOptions() models.ProjectOptions {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Printf(
		"\n%s\n%s\n\n",
		styles.Bold(styles.Accent("GoStart")),
		styles.Muted("Create a new Go project"),
	)

	return models.ProjectOptions{
		ProjectName: prompt.Ask(
			scanner,
			"Project name "+styles.Muted("[my-go-project]"),
		),
		ModuleName: prompt.Ask(
			scanner,
			"Go module "+styles.Muted("[example.com/go-module]"),
		),
		ProjectType: prompt.Choice(scanner, &prompt.ChoiceOptions{
			Title:          "Project type",
			Prompt:         "Select",
			Options:        []string{"basic", "http", "cli"},
			Default:        "basic",
			InvalidMessage: "Unknown project type. Choose basic, http, or cli",
		}),
		InitGit: prompt.Confirm(scanner, "Initialize Git repository?"),
	}
}

func RunCreate() error {
	options := collectProjectOptions()

	return generate.GenerateProject(options)
}
