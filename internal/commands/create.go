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
		ProjectName: prompt.Ask(scanner, "Название проекта"),
	}
}

func RunCreate() error {
	options := collectProjectOptions()

	return generate.GenerateProject(options)
}
