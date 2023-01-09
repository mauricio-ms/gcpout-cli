/*
Copyright © 2022 Maurício Mussatto Scopel <ms.mauricio93@gmail.com>

*/
package cmd

import (
	"fmt"
	"log"
	"errors"
	"github.com/fatih/color"
	"sort"
	"strings"
	"path/filepath"
	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/core"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
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
	      configFile, err := ReadConfigFile()
	      if err != nil {
	      	 Error("X Run the init command to configure the CLI\n")
		 return
	      }

	      projectsPath, err := ProjectsPath()
	      if err != nil {
	      	 log.Fatal(err.Error())
		 return
	      }
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

	      branches, err := repositoryClone.Branches()
	      if err != nil {
	      	 log.Fatal(err.Error())
		 return
	      }

	      currentBranch, err := repositoryClone.CurrentBranch()
	      if err != nil {
	      	 log.Fatal(err.Error())
		 return
	      }

	      var sourceBranch string
	      survey.AskOne(
		&survey.Select{
			Message: "Source Branch:",
			Options: branches,
			Default: currentBranch,
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
	      remoteBranches, err := repositoryClone.RemoteBranches()
	      if err != nil {
	      	 log.Fatal(err.Error())
		 return
	      }
	      sourceBranchIndex, _ := sort.Find(len(remoteBranches), func(i int) int {
	      			 return strings.Compare(sourceBranch, remoteBranches[i])
	      })
	      targetBranches := make([]string, 0)
	      targetBranches = append(targetBranches, remoteBranches[:sourceBranchIndex]...)
	      targetBranches = append(targetBranches, remoteBranches[sourceBranchIndex+1:]...)

	      var targetBranch string
	      survey.AskOne(
		&survey.Select{
			Message: "Target Branch:",
			Options: targetBranches,
		}, &targetBranch)


	      var jiraIssueId string
	      survey.AskOne(
		&survey.Input{
			Message: "Jira Issue Id:",
			Help: "For a project ST and an issue 123, type ST-123",
		}, &jiraIssueId, survey.WithValidator(JiraIssueValidator()))

	      jira := NewJira(configFile.jiraServer)
	      pullRequestTemplate := NewPullRequestTemplate(jira.LinkForIssue(jiraIssueId))

	      var qs = []*survey.Question{
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
						Options: pullRequestTemplate.TypeOfChanges,
					},
		  },
		  {
			Name:		"checklist",
			Prompt:		&survey.MultiSelect{
						Message: "Checklist:",
						Options: pullRequestTemplate.ChecklistQuestions,
					},
		  },
	      }

	      answers := struct {
		      Description	string
		      TypeOfChange	int
		      Checklist		[]core.OptionAnswer
	      }{}

	      err = survey.Ask(qs, &answers)
	      if err != nil {
	      	 log.Fatal(err.Error())
		 return
	      }

	      pullRequestTemplate.Description = answers.Description
	      pullRequestTemplate.TypeOfChange = answers.TypeOfChange
	      pullRequestTemplate.Checklist = make([]int, len(answers.Checklist))
	      for i, checklistAnswer := range answers.Checklist {
	      	  pullRequestTemplate.Checklist[i] = checklistAnswer.Index
	      }

	      link, err := repositoryClone.OpenPullRequest(jiraIssueId, pullRequestTemplate, sourceBranch, targetBranch)
	      if err != nil {
	      	 Errorf("X %s\n", err.Error())
	      } else {
	      	 Successf("PR opened: %s\n", link)
	      }
	},
}

func Successf(message string, args ...any) {
     color.New(color.FgGreen).Printf(message, args...)
}

func Error(message string) {
     color.New(color.FgRed).Print(message)
}

func Errorf(message string, args ...any) {
     color.New(color.FgRed).Printf(message, args...)
}

func RemoteBranchValidator(rc RepositoryClone) survey.Validator {
     return func (val interface{}) error {
     	    answer, ok := val.(core.OptionAnswer)
	    if !ok {
	       return fmt.Errorf("Internal error")
	    }

	    hasRemoteBranch, err := rc.HasRemoteBranch(answer.Value)
	    if err != nil {
	       return err
	    }
	       
	    if !hasRemoteBranch {
	       return fmt.Errorf("%s is only local", answer.Value)
	    }
	    return nil
     }
}

func JiraIssueValidator() survey.Validator {
     return func (val interface{}) error {
     	    if value, ok := val.(string) ; !ok || !ValidJiraIssue(value) {
	       return fmt.Errorf("%s doesn't follow the Jira Issue pattern <PROJECT_ID-ISSUE_ID>", value)
	    } 
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

func ProjectsPath() (string, error) {
     var inner func(relativePath string) (string, error)
     inner = func (relativePath string) (string, error) {
     	     	  if _, err := os.Stat(relativePath + ".git"); err == nil {
     	       	     return inner(relativePath + "../")
     	     	  }
     	     	  return RunCommand("readlink", "-f", relativePath)
     	     }
     path, err := inner("")
     if err != nil {
     	return "", err
     }
     return path + "/", nil
}

func RunCommand(name string, arg ...string) (string, error) {
      command := exec.Command(name, arg...)

      output, err := command.CombinedOutput()
      if err != nil {
      	 errValue := string(output)
      	 endJsonIdx := strings.LastIndex(errValue, "}")
  	 return "", errors.New(errValue[0:endJsonIdx+1])
      }
      return strings.TrimSuffix(string(output), "\n"), nil
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