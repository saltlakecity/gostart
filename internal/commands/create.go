package commands

import (
	"bufio"
	"os"

	"github.com/saltlakecity/gostart/internal/generate"
	"github.com/saltlakecity/gostart/internal/models"
	"github.com/saltlakecity/gostart/internal/prompt"
)

func collectProjectOptions() models.ProjectOptions {
	scanner := bufio.NewScanner(os.Stdin)

	return models.ProjectOptions{
		ProjectName: prompt.Ask(scanner, "Enter project name (or leave empty for default project name)"),
		ModuleName:  prompt.Ask(scanner, "Enter Go-module name"),
		InitGit:     prompt.Confirm(scanner, "Init Git repository?"),
	}
}

func RunCreate() error {
	options := collectProjectOptions()

	return generate.GenerateProject(options)
}
