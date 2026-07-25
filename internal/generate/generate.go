package generate

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/saltlakecity/gostart/internal/models"
)

func initGitRepo(path string) error {

	if _, err := exec.LookPath("git"); err != nil {
		return fmt.Errorf("git was not found: Install git or add git to PATH")
	}

	cmd := exec.Command("git", "init")

	cmd.Dir = path

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Can't resolve <git init>: %w", err)
	}

	return nil
}
func GenerateProject(opts models.ProjectOptions) error {
	name := opts.ProjectName
	git := opts.InitGit
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
	if git {
		err := initGitRepo(filepath.Join(name))
		if err != nil {
			return err
		}
	}
	return nil

}
