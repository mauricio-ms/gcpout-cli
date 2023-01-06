package cmd

import (
       "fmt"
       "regexp"
)

type Jira struct {
     server	 string
}

func NewJira(server string) Jira {
     return Jira {
     	    server: server,
     }
}

func (this Jira) LinkForIssue(issueId string) string {
     return fmt.Sprintf("%s/browse/%s", this.server, issueId)
}

func ValidJiraHost(host string) bool {
     match, err := regexp.MatchString(`^https://\w+.atlassian.net$`, host)
     if err != nil {
     	return false
     } else {
        return match
     }
}

func ValidJiraIssue(issueId string) bool {
     match, err := regexp.MatchString(`^\w+-\d+$`, issueId)
     if err != nil {
     	return false
     } else {
        return match
     }
}