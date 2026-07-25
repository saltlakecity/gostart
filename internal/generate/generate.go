package generate

import (
	"os"
	"path/filepath"

	"github.com/saltlakecity/gostart/internal/models"
)

func GenerateProject(opts models.ProjectOptions) error {
	name := opts.ProjectName

	if name == "" {
		name = "my-go-project"
	}

	dirs := []string{
		filepath.Join(name),
	}

	for _, dir := range dirs {
		err := os.MkdirAll(dir, 0644)

		if err != nil {
			return err
		}
	}

	return nil

}
