package generate

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/saltlakecity/gostart/internal/models"
	"github.com/saltlakecity/gostart/internal/styles"
	projectTemplates "github.com/saltlakecity/gostart/internal/templates"
)

func createGitIgnore(projectRoot string) error {
	content := `.env
bin/
dist/
*.exe
`

	path := filepath.Join(projectRoot, ".gitignore")

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("create .gitignore: %w", err)
	}

	return nil
}

func initGoModule(projectRoot, moduleName string) error {
	if moduleName == "" {
		moduleName = "example.com/go-module"
	}

	cmd := exec.Command("go", "mod", "init", moduleName)
	cmd.Dir = projectRoot

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf(
			"create Go module: %w: %s",
			err,
			strings.TrimSpace(string(output)),
		)
	}

	return nil
}

func initGitRepo(projectRoot string) error {
	if _, err := exec.LookPath("git"); err != nil {
		return fmt.Errorf(
			"git was not found: install Git or add it to PATH",
		)
	}

	cmd := exec.Command("git", "init", "--quiet")
	cmd.Dir = projectRoot

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf(
			"execute git init: %w: %s",
			err,
			strings.TrimSpace(string(output)),
		)
	}

	return nil
}

func checkDestinationAvailable(projectRoot string) error {
	_, err := os.Stat(projectRoot)

	switch {
	case err == nil:
		return fmt.Errorf(
			"destination folder %q already exists",
			projectRoot,
		)
	case errors.Is(err, os.ErrNotExist):
		return nil
	default:
		return fmt.Errorf(
			"inspect destination folder %q: %w",
			projectRoot,
			err,
		)
	}

}

func GenerateProject(
	opts models.ProjectOptions,
) (resultErr error) {
	name := opts.ProjectName
	moduleName := opts.ModuleName
	projectType := opts.ProjectType

	if name == "" {
		name = "my-go-project"
	}

	if projectType == "" {
		projectType = "basic"
	}

	opts.ProjectName = name
	opts.ProjectType = projectType

	projectRoot := name

	if err := checkDestinationAvailable(projectRoot); err != nil {
		return err
	}

	templateFS, err := fs.Sub(
		projectTemplates.Files,
		projectType,
	)
	if err != nil {
		return fmt.Errorf(
			"select project template %q: %w",
			projectType,
			err,
		)
	}

	parentDirectory := filepath.Dir(projectRoot)
	temporaryPattern := "." + filepath.Base(projectRoot) + "-*"

	stagingRoot, err := os.MkdirTemp(
		parentDirectory,
		temporaryPattern,
	)
	if err != nil {
		return fmt.Errorf(
			"create temporary project directory: %w",
			err,
		)
	}

	removeStaging := true

	defer func() {
		if !removeStaging {
			return
		}

		if cleanupErr := os.RemoveAll(stagingRoot); cleanupErr != nil {
			resultErr = errors.Join(
				resultErr,
				fmt.Errorf(
					"remove temporary project directory %s: %w",
					stagingRoot,
					cleanupErr,
				),
			)
		}
	}()

	fmt.Printf("\n%s\n", styles.Bold("Creating project"))

	if err := GenerateFromTemplate(
		stagingRoot,
		templateFS,
		opts,
	); err != nil {
		return fmt.Errorf(
			"generate project from template: %w",
			err,
		)
	}

	if err := initGoModule(stagingRoot, moduleName); err != nil {
		return err
	}
	printCompleted("Go module initialized")

	if err := createGitIgnore(stagingRoot); err != nil {
		return err
	}
	printCompleted(".gitignore created")

	if opts.InitGit {
		if err := initGitRepo(stagingRoot); err != nil {
			return err
		}
		printCompleted("Git repository initialized")
	}

	if err := os.Rename(stagingRoot, projectRoot); err != nil {
		return fmt.Errorf(
			"move completed project to %s: %w",
			projectRoot,
			err,
		)
	}

	removeStaging = false

	fmt.Printf(
		"\n%s %s\n\n%s\n  %s\n  %s\n",
		styles.Success(styles.SuccessIcon),
		styles.Bold(
			fmt.Sprintf(
				"Project %q successfully created",
				name,
			),
		),
		styles.Muted("Next steps:"),
		styles.Accent("cd "+name),
		styles.Accent("go run ."),
	)

	return nil
}

func printCompleted(message string) {
	fmt.Printf(
		"  %s %s\n",
		styles.Success(styles.SuccessIcon),
		message,
	)
}

func GenerateFromTemplate(
	projectRoot string,
	sourceFS fs.FS,
	data any,
) error {
	return fs.WalkDir(
		sourceFS,
		".",
		func(
			sourcePath string,
			entry fs.DirEntry,
			walkErr error,
		) error {
			if walkErr != nil {
				return fmt.Errorf(
					"walk through %s: %w",
					sourcePath,
					walkErr,
				)
			}

			destinationPath := filepath.Join(
				projectRoot,
				filepath.FromSlash(sourcePath),
			)

			if entry.IsDir() {
				if err := os.MkdirAll(destinationPath, 0755); err != nil {
					return fmt.Errorf(
						"create directory %s: %w",
						destinationPath,
						err,
					)
				}

				fmt.Printf(
					"  %s %s\n",
					styles.Muted("mkdir "),
					styles.Muted(destinationPath),
				)
				return nil
			}

			if err := os.MkdirAll(
				filepath.Dir(destinationPath),
				0755,
			); err != nil {
				return fmt.Errorf(
					"create parent directory for %s: %w",
					destinationPath,
					err,
				)
			}

			if strings.HasSuffix(sourcePath, ".tmpl") {
				destinationPath = strings.TrimSuffix(
					destinationPath,
					".tmpl",
				)

				if err := renderTemplate(
					sourceFS,
					sourcePath,
					destinationPath,
					data,
				); err != nil {
					return err
				}

				fmt.Printf(
					"  %s %s\n",
					styles.Success("create"),
					destinationPath,
				)
				return nil
			}

			if err := copyFile(
				sourceFS,
				sourcePath,
				destinationPath,
			); err != nil {
				return err
			}

			fmt.Printf(
				"  %s %s\n",
				styles.Success("create"),
				destinationPath,
			)
			return nil
		},
	)
}

func renderTemplate(
	sourceFS fs.FS,
	sourcePath string,
	destinationPath string,
	data any,
) error {
	content, err := fs.ReadFile(sourceFS, sourcePath)
	if err != nil {
		return fmt.Errorf(
			"read template %s: %w",
			sourcePath,
			err,
		)
	}

	tmpl, err := template.New(sourcePath).Parse(string(content))
	if err != nil {
		return fmt.Errorf(
			"parse template %s: %w",
			sourcePath,
			err,
		)
	}

	outputFile, err := os.Create(destinationPath)
	if err != nil {
		return fmt.Errorf(
			"create output file %s: %w",
			destinationPath,
			err,
		)
	}

	if err := tmpl.Execute(outputFile, data); err != nil {
		outputFile.Close()
		_ = os.Remove(destinationPath)

		return fmt.Errorf(
			"execute template %s: %w",
			sourcePath,
			err,
		)
	}

	if err := outputFile.Close(); err != nil {
		return fmt.Errorf(
			"close output file %s: %w",
			destinationPath,
			err,
		)
	}

	return nil
}

func copyFile(
	sourceFS fs.FS,
	sourcePath string,
	destinationPath string,
) error {
	sourceFile, err := sourceFS.Open(sourcePath)
	if err != nil {
		return fmt.Errorf(
			"open source file %s: %w",
			sourcePath,
			err,
		)
	}
	defer sourceFile.Close()

	destinationFile, err := os.OpenFile(
		destinationPath,
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
		0644,
	)
	if err != nil {
		return fmt.Errorf(
			"create destination file %s: %w",
			destinationPath,
			err,
		)
	}

	if _, err := io.Copy(destinationFile, sourceFile); err != nil {
		destinationFile.Close()

		return fmt.Errorf(
			"copy %s to %s: %w",
			sourcePath,
			destinationPath,
			err,
		)
	}

	if err := destinationFile.Close(); err != nil {
		return fmt.Errorf(
			"close destination file %s: %w",
			destinationPath,
			err,
		)
	}

	return nil
}
