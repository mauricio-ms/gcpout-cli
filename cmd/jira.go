package cmd

import (
       "fmt"
)

// TODO read values from a file like /home/user/.config/.gcpout/.config.properties
// 	create the file with a init command

type Jira struct {
     server	 string
}

func NewJira() Jira {
     return Jira {
     	    server: "https://localhost", // TODO read from the file
     }
}

func (this Jira) LinkForIssue(issueId string) string {
     return fmt.Sprintf("%s/browse/%s", this.server, issueId)
}