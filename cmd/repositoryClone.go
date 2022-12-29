package cmd

import (
       "regexp"
       "encoding/json"
       "fmt"
       "strings"
       "sort"
)

type RepositoryClone struct {
     ParentPath	       string
     Name	       string
}

// TODO Create New method to precompute the methods into fields

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

func (rc RepositoryClone) Branches() []string {
     branchesMap := map[string]struct{}{}

     for _, branch := range rc.LocalBranches() {
     	 branchesMap[branch] = struct{}{}
     }
     for _, branch := range rc.RemoteBranches() {
     	 branchesMap[branch] = struct{}{}
     }

     branches := make([]string, 0, len(branchesMap))
     for branch := range branchesMap {
     	 branches = append(branches, branch)
     }

     sort.Strings(branches)

     return branches
}

func (rc RepositoryClone) LocalBranches() []string {
     listBranches := runCommand("git", "-C", rc.Path(), "branch", "--list")
     var localBranches []string
     localBranches = strings.Split(listBranches, "\n")

     for i, branch := range localBranches {
     	 branch = strings.TrimSpace(branch)
	 branch = strings.TrimPrefix(branch, "* ")
     	 localBranches[i] = branch
     }
     
     return localBranches
}

func (rc RepositoryClone) HasRemoteBranch(branch string) bool {
     remoteBranches := rc.RemoteBranches()
     _, found := sort.Find(len(remoteBranches), func(i int) int {
     	return strings.Compare(branch, remoteBranches[i])
     })

     return found
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