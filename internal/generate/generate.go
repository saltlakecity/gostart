package generate

import (
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

func GenerateProject(opts models.ProjectOptions) error {
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
	templateRoot := filepath.Join("internal", "templates")

	fmt.Printf("\n%s\n", styles.Bold("Creating project"))

	if err := GenerateFromTemplate(
		projectRoot,
		templateRoot,
		projectType,
		opts,
	); err != nil {
		return fmt.Errorf("generate project from template: %w", err)
	}

	if err := initGoModule(projectRoot, moduleName); err != nil {
		return err
	}
	printCompleted("Go module initialized")

	if err := createGitIgnore(projectRoot); err != nil {
		return err
	}
	printCompleted(".gitignore created")

	if opts.InitGit {
		if err := initGitRepo(projectRoot); err != nil {
			return err
		}
		printCompleted("Git repository initialized")
	}

	fmt.Printf(
		"\n%s %s\n\n%s\n  %s\n  %s\n",
		styles.Success(styles.SuccessIcon),
		styles.Bold(fmt.Sprintf("Project %q successfully created", name)),
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
	templateRoot string,
	projectType string,
	data any,
) error {
	sourceRoot := filepath.Join(templateRoot, projectType)

	return filepath.WalkDir(
		sourceRoot,
		func(path string, entry fs.DirEntry, walkErr error) error {
			if walkErr != nil {
				return fmt.Errorf("walk through %s: %w", path, walkErr)
			}

			relativePath, err := filepath.Rel(sourceRoot, path)
			if err != nil {
				return fmt.Errorf(
					"calculate relative path for %s: %w",
					path,
					err,
				)
			}

			destinationPath := filepath.Join(projectRoot, relativePath)

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

			if strings.HasSuffix(path, ".tmpl") {
				destinationPath = strings.TrimSuffix(
					destinationPath,
					".tmpl",
				)

				if err := renderTemplate(
					path,
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

			if err := copyFile(path, destinationPath); err != nil {
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
	sourcePath string,
	destinationPath string,
	data any,
) error {
	tmpl, err := template.ParseFiles(sourcePath)
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
		os.Remove(destinationPath)

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

func copyFile(sourcePath, destinationPath string) error {
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf(
			"open source file %s: %w",
			sourcePath,
			err,
		)
	}
	defer sourceFile.Close()

	sourceInfo, err := sourceFile.Stat()
	if err != nil {
		return fmt.Errorf(
			"read file information for %s: %w",
			sourcePath,
			err,
		)
	}

	destinationFile, err := os.OpenFile(
		destinationPath,
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
		sourceInfo.Mode(),
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
