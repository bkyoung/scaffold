/*
Copyright © 2020 Brandon Young <bkyoung@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/bkyoung/scaffold/internal/git"
	"github.com/bkyoung/scaffold/internal/scaffold"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	project    scaffold.Project
	repo       git.GithubRepository
	createRepo bool
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize project",
	Long:  `Initialize a new go project`,
	Args:  cobra.ExactArgs(1),
	PreRun: func(cmd *cobra.Command, args []string) {
		// Process config, environment variables, and flags
		viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
		viper.AutomaticEnv()

		// Populate as much project info as we can from viper
		err := viper.Unmarshal(&project)
		if err != nil {
			fmt.Printf("could not retrieve supplied project settings: %s\n", err)
			os.Exit(1)
		}

		// Additional arg processing
		project.Configure(scaffold.Name(args[0]), scaffold.ProjectDir(project.ProjectDir))

		if project.CreateRepo {
			if err := viper.Unmarshal(&repo); err != nil {
				fmt.Printf("error unpacking repo info: %s\n", err)
				os.Exit(1)
			}
			repo.Name = project.Name
			project.Repo = repo
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Create the remote/origin repo, if requested
		if createRepo {
			// Do the thing
			err := project.Repo.Create()
			if err != nil {
				fmt.Printf("error creating repository: %s\n", err)
				os.Exit(1)
			}

			// Set the go module's name according to the repo URL
			// TODO: better algo for this, so we can use it with ANY conn type
			if url, err := project.Repo.URL(); err == nil && len(url) > 9{
				project.GoModuleName = url[8:]
			}
		}

		// Clone the repo
		if err := project.Repo.Clone(os.Stdout); err != nil {
			fmt.Printf("error cloning repository: %s\n", err)
			os.Exit(1)
		}

		// If the go module's name is still not set, just set it to the project name
		if project.GoModuleName == "" {
			project.GoModuleName = project.Name
		}

		// Create the local project structure
		err := scaffold.Create(project)
		if err != nil {
			fmt.Printf("error creating project directory: %s\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().Bool("create-repo", false, "Create a new remote/origin for this project (default is false)")
	initCmd.Flags().Bool("disable-modules", false, "Disable go modules for this project (default is false)")
	initCmd.Flags().String("github-auth-token", "", "Github (Enterprise) Personal Auth Token for creating remote/origin (Ignored unless --create-repo used)")
	initCmd.Flags().String("github-enterprise-server-url", "", "Base URL of Github Enterprise server if needed (ex: https://ghe.example.com) (Ignored unless --create-repo used)")
	initCmd.Flags().String("github-organization", "", "Org owner of remote/origin if not the user (Ignored unless --create-repo used)")
	initCmd.Flags().String("github-repo-name", "", "Repository name to create (default is the project name) (Ignored unless --create-repo used)")
	initCmd.Flags().String("license", "mit", "Project license (Ignored unless --create-repo used)")
	initCmd.Flags().String("module-name", "", "Name of the go module for this project (default is the project name)")
	initCmd.Flags().Bool("private", false, "Make new repo private (default is not private) (Ignored unless --create-repo used)")
	initCmd.Flags().String("project-dir", "", "Directory in which to initialize new project (default is project name)")
	err := viper.BindPFlags(initCmd.Flags())
	if err != nil {
		fmt.Printf("error reading flags: %s\n", err)
		os.Exit(1)
	}
}
