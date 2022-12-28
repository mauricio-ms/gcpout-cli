package cmd

import (
       "regexp"
       "encoding/json"
       "fmt"
)

type RepositoryClone struct {
     ParentPath	       string
     Name	       string
}

func (rc RepositoryClone) Path() string {
     return rc.ParentPath + "/" + rc.Name
}

func (rc RepositoryClone) CurrentBranch() string {
     return runCommand("git", "-C", rc.Path(), "branch", "--show-current")
}

func (rc RepositoryClone) RepoOwner() string {
     originUrl := runCommand("git", "-C", rc.Path(), "remote", "get-url", "origin")
     ownerRegex := regexp.MustCompile(`^[^:]+:([^/]+).*$`)
     
     return ownerRegex.FindStringSubmatch(originUrl)[1]
}

func (rc RepositoryClone) LocalBranches() []string {
     // TODO implement it
     return make([]string, 0)
}

func (rc RepositoryClone) RemoteBranches() []string {
     endpoint := rc.ApiResourcePrefix() + "/branches"
     response := runCommand("gh", "api", "-H", "Accept:application/vnd.github+json", endpoint)

     var branches []map[string]string
     json.Unmarshal([]byte(response), &branches)

     var branchesNames = make([]string, len(branches))
     
     for i, branch := range branches {
     	 branchesNames[i] = branch["name"]
     }
     
     return branchesNames
}

func (rc RepositoryClone) ApiResourcePrefix() string {
     return fmt.Sprintf("/repos/%s/%s", rc.RepoOwner(), rc.Name)
}