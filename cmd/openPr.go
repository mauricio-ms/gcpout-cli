/*
Copyright © 2022 Maurício Mussatto Scopel <ms.mauricio93@gmail.com>

*/
package cmd

import (
	"fmt"
	"log"
	"github.com/fatih/color"
	"sort"
	"strings"
	"path/filepath"
	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/core"
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

	      var sourceBranch string
	      survey.AskOne(
		&survey.Select{
			Message: "Source Branch:",
			Options: repositoryClone.Branches(),
			Default: repositoryClone.CurrentBranch(),
		}, &sourceBranch, survey.WithValidator(RemoteBranchValidator(repositoryClone)))

	      /**
	      TODO use it when new questions is used to push the branch to user
	      if (!repositoryClone.HasRemoteBranch(sourceBranch)) {
		 // TODO add code to ask if user can push his local branch
		 // example: survey.AskOne(&survey.Confirm{Message: ""})
		 Errorf("X The branch %s is only local, please check if everithing is commited and push it before opening the PRs\n", sourceBranch)
		 return
	      }
	      **/

	      // TODO create new struct
	      remoteBranches := repositoryClone.RemoteBranches()
	      sourceBranchIndex, _ := sort.Find(len(remoteBranches), func(i int) int {
	      			 return strings.Compare(sourceBranch, remoteBranches[i])
	      })
	      targetBranches := make([]string, 0)
	      targetBranches = append(targetBranches, remoteBranches[:sourceBranchIndex]...)
	      targetBranches = append(targetBranches, remoteBranches[sourceBranchIndex+1:]...)

	      var qs = []*survey.Question{
		  {
			Name:		"targetBranch",
			Prompt:		&survey.Select{
						Message: "Target Branch:",
						Options: targetBranches,
					},
		  },
	      	  {
			Name:		"title",
			Prompt:		&survey.Input{Message: "PR Title:"},
			Validate: 	survey.Required,
		  },
		  {
			Name:		"description",
			Prompt:		&survey.Input{
						Message: "Description:",
						Help: "Please include a summary of the change and which issue is fixed. Please also include relevant motivation and context. List any dependencies that are required for this change.",
					},
			Validate: 	survey.Required,
		  },
		  {
			Name:		"typeOfChange",
			Prompt:		&survey.Select{
						Message: "Type of change:",
						Options: []string{
							 "Chore (no production code changed. eg: doc, style, tests)Chore (no production code changed. eg: doc, style, tests)",
							 "Bug fix (non-breaking change which fixes an issue)",
							 "New feature (non-breaking change which adds functionality)",
							 "Breaking change (fix or feature that would cause existing functionality to not work as expected)",
						},
					},
		  },
		  {
			Name:		"checklist",
			Prompt:		&survey.MultiSelect{
						Message: "Checklist:",
						Options: []string{
							 "I have commented on my code, particularly in hard-to-understand areas",
							 "I have made corresponding changes to the documentation",
							 "The required manual tests were executed and no issue was found",
						},
					},
		  },
	      }

	      answers := struct {
		      TargetBranch	string
	      	      Title		string
		      Description	string
		      TypeOfChange	int
		      Checklist		[]core.OptionAnswer
	      }{}

	      err := survey.Ask(qs, &answers)
	      if err != nil {
	      	 log.Fatal(err.Error())
		 return
	      }

	      fmt.Println(answers)

	      return

	      endpoint := fmt.Sprintf("/repos/%s/%s/pulls", repositoryClone.RepoOwner(), project)

	      openPrCommand := runCommand("gh", "api",
	      		    "--method", "POST", "-H", "Accept:application/vnd.github+json", endpoint,
			    "-f", fmt.Sprintf("title=%s", answers.Title),
			    "-f", fmt.Sprintf("body=%s", answers.Description),
			    "-f", fmt.Sprintf("head=%s", sourceBranch),
			    "-f", fmt.Sprintf("base=%s", answers.TargetBranch))
	      var pr map[string]string
	      json.Unmarshal([]byte(openPrCommand), &pr)
	      fmt.Println("PR opened: " + pr["html_url"])
	},
}

func Errorf(message string, args ...string) {
     color.New(color.FgRed).Printf(message, args)
}

func RemoteBranchValidator(rc RepositoryClone) survey.Validator {
     return func (val interface{}) error {
     	    if answer, ok := val.(core.OptionAnswer) ; !ok || !rc.HasRemoteBranch(answer.Value) {
	       return fmt.Errorf("%s is only local", answer.Value)
	    }
	    fmt.Println("Ok")
	    return nil
     }
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