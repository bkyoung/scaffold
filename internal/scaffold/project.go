package scaffold

import (
	"fmt"
	"github.com/bkyoung/scaffold/internal/repository"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/pkg/errors"
)

type Project struct {
	Name            string
	Repo 			repository.SCMRepository
	CreateRepo 		bool   `mapstructure:"create-repo"`
	ProjectDir      string `mapstructure:"project-dir"`
	DisableModules  bool   `mapstructure:"disable-modules"`
	GoModuleName    string `mapstructure:"module-name"`
}

type option func(*Project)

func (p *Project) Configure(opts ...option) {
	for _, opt := range opts {
		opt(p)
	}
}

func Name(n string) option {
	return func(p *Project) {
		p.Name = n
	}
}

func CreateRepo(cr bool) option {
	return func(p *Project) {
		p.CreateRepo = cr
	}
}

func ProjectDir(dir string) option {
	return func(p *Project) {
		// Ensure the ProjectDir is an absolute path
		pd, _ := filepath.Abs(dir)
		p.ProjectDir = pd
	}
}

func DisableModules(dm bool) option {
	return func(p *Project) {
		p.DisableModules = dm
	}
}

func GoModuleName(mod string) option {
	return func(p *Project) {
		// Ensure we trim off leading 'https://' or 'http://'
		name := strings.TrimPrefix(mod, "https://")
		name = strings.TrimPrefix(mod, "http://")
		p.GoModuleName = name
	}
}

func Create(p Project) error {
	// Save our starting place
	curDir, err := os.Getwd()
	if err != nil {
		return errors.Wrap(err,"failed to determine current working directory")
	}

	if curDir != p.ProjectDir {
		if _, err := os.Stat(p.ProjectDir); os.IsNotExist(err) {
			err = os.MkdirAll(p.ProjectDir, 0755)
			if err != nil {
				return errors.Wrap(err, "failed to create project directory")
			}
		}

		err = os.Chdir(p.ProjectDir)
		if err != nil {
			return errors.Wrap(err, "failed to change into project directory")
		}
	}

	// Render a README
	t, err := template.New("readme").Parse("# {{ .Name}}\nThis is the home of {{ .Name }}\n")
	if err != nil {
		return errors.Wrap(err, "error parsing readme template")
	}
	r, err := os.Create("README.md");if err != nil {
		return errors.Wrap(err, "error creating Readme.md file")
	}
	err = t.Execute(r, p)
	if err != nil {
		return errors.Wrap(err, "error rendering readme template")
	}

	// Initialize go modules
	if !p.DisableModules {
		modInit := exec.Command("go", "mod", "init", p.GoModuleName)
		output, err := modInit.CombinedOutput();if err != nil {
			fmt.Printf("%s\n", string(output))
			errors.Wrap(err, "error while initializing go modules in project directory")
		}
	}

	// Return to our starting place
	err = os.Chdir(curDir);if err != nil {
		return errors.Wrap(err, "failed to change back to initial directory from project directory")
	}

	return nil
}
