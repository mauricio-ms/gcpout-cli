/*
Copyright © 2022 Maurício Mussatto Scopel <ms.mauricio93@gmail.com>

*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/AlecAivazis/survey/v2"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		var jiraHost string
	        survey.AskOne(
			&survey.Input{
				Message: "Jira Host:",
				Help: "https://<your_server>.atlassian.net",
			}, &jiraHost, survey.WithValidator(JiraHostValidator()))
		
		InitConfigFile(jiraHost).Persist()
	},
}

func JiraHostValidator() survey.Validator {
     return func (val interface{}) error {
     	    if value, ok := val.(string) ; !ok || !ValidJiraHost(value) {
	       return fmt.Errorf("%s doesn't follow the Jira Host https://<your_server>.atlassian.net", value)
	    } 
     	    return nil
     }
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
