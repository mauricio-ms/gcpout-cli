/*
Copyright © 2022 Maurício Mussatto Scopel <ms.mauricio93@gmail.com>

*/
package cmd

import (
	"fmt"
	"log"
	"strings"
	"path/filepath"
	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"encoding/json"
)

// openPrCmd represents the openPr command
var openPrCmd = &cobra.Command{
	Use:   "openPr",
	Short: "Open a Pull Request",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
   	//Args: cobra.MinimumNArgs(1),
   	Run: func(cmd *cobra.Command, args []string) {
	      projectsPath := projectsPath()
	      projects := projects(projectsPath + "*")

	      var project string
	      survey.AskOne(
		&survey.Select{
			Message: "Project:",
			Options: projects,
		}, &project, survey.WithValidator(survey.Required))

	      repositoryClone := RepositoryClone {
	      		      ParentPath: projectsPath,
			      Name: project,
	      }

	      var qs = []*survey.Question{
		  {
			Name:		"sourceBranch",
			Prompt:		&survey.Input{
						Message: "Source Branch:",
						Default: repositoryClone.CurrentBranch(),
					},
			Validate:	survey.Required,
			// TODO add options listing all branches for the project
			// TODO add validation: should be a valid remote branch
		  },
		  {
			Name:		"targetBranch",
			Prompt:		&survey.Input{
						Message: "Target Branch:",
					},
			Validate:	survey.Required,
			// TODO add options listing all branches (using LocalBranches + RemoteBranches methods) for the project but the choosed source branch
			// TODO add validation: should be a valid remote branch
		  },
	      	  {
			Name:		"title",
			Prompt:		&survey.Input{Message: "PR Title:"},
			Validate: 	survey.Required,
		  },
		  {
			Name:		"description",
			Prompt:		&survey.Input{Message: "PR Description:"},
			Validate: 	survey.Required,
		  },
	      }

	      answers := struct {
	      	      Project		string
		      SourceBranch	string
		      TargetBranch	string
	      	      Title		string
		      Description	string
	      }{}

	      err := survey.Ask(qs, &answers)
	      if err != nil {
	      	 log.Fatal(err.Error())
		 return
	      }

	      fmt.Println(repositoryClone.RepoOwner())

	      return

	      endpoint := fmt.Sprintf("/repos/%s/%s/pulls", repositoryClone.RepoOwner(), answers.Project)

	      openPrCommand := runCommand("gh", "api",
	      		    "--method", "POST", "-H", "Accept:application/vnd.github+json", endpoint,
			    "-f", fmt.Sprintf("title=%s", answers.Title),
			    "-f", fmt.Sprintf("body=%s", answers.Description),
			    "-f", fmt.Sprintf("head=%s", answers.SourceBranch),
			    "-f", fmt.Sprintf("base=%s", answers.TargetBranch))
	      var pr map[string]string
	      json.Unmarshal([]byte(openPrCommand), &pr)
	      fmt.Println("PR opened: " + pr["html_url"])
	},
}

func projects(projectsPath string) []string {
     var projectsPaths, _ = filepath.Glob(projectsPath)
     var projects = make([]string, len(projectsPaths))
     for i, v := range projectsPaths {
     	 hierarchy := strings.Split(v, "/")
	 projects[i] = hierarchy[len(hierarchy)-1]
     }
     return projects
}

func projectsPath() string {
     var inner func(relativePath string) string
     inner = func (relativePath string) string {
     	     	  if _, err := os.Stat(relativePath + ".git"); err == nil {
     	       	     return inner(relativePath + "../")
     	     	  }
     	     	  return runCommand("readlink", "-f", relativePath)
     	     }
     return inner("") + "/"
}

// TODO add behavior to suppress errors only when needed
func runCommand(name string, arg ...string) string {
      command := exec.Command(name, arg...)

      output, err := command.CombinedOutput()
      if err != nil {
  	 return ""
      }
      return strings.TrimSuffix(string(output), "\n")
}

func init() {
	rootCmd.AddCommand(openPrCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// openPrCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// openPrCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}