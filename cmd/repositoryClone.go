package cmd

import (
       "regexp"
       "encoding/json"
       "fmt"
       "strings"
       "sort"
       "errors"
)

type RepositoryClone struct {
     ParentPath	       string
     Name	       string
     Path	       string
     RepoOwner	       string
     ApiResourcePrefix string
     CurrentBranch     string
     LocalBranches     []string
     RemoteBranches    []string
     Branches	       []string
}

func GetRepositoryClone(parentPath string, name string) (*RepositoryClone, error) {
     rc := &RepositoryClone {
     	ParentPath: parentPath,
	Name: name,
	Path: parentPath + "/" + name,
     }

     var err error
     rc.RepoOwner, err = _RepoOwner(rc.Path)
     if err != nil {
     	return nil, err
     }

     rc.ApiResourcePrefix = fmt.Sprintf("/repos/%s/%s", rc.RepoOwner, rc.Name)

     rc.CurrentBranch, err = _CurrentBranch(rc.Path)
     if err != nil {
     	return nil, err
     }
     
     rc.LocalBranches, err = _LocalBranches(rc.Path)
     if err != nil {
     	return nil, err
     }

     rc.RemoteBranches, err = _RemoteBranches(rc.ApiResourcePrefix)
     if err != nil {
     	return nil, err
     }

     rc.Branches = _Branches(rc.LocalBranches, rc.RemoteBranches)

     return rc, nil
}

func _RepoOwner(path string) (string, error) {
     originUrl, err := RunCommand("git", "-C", path, "remote", "get-url", "origin")
     if err != nil {
     	return "", err
     }
     
     ownerRegex := regexp.MustCompile(`^[^:]+:([^/]+).*$`)
     return ownerRegex.FindStringSubmatch(originUrl)[1], nil
}

func _LocalBranches(path string) ([]string, error) {
     listBranches, err := RunCommand("git", "-C", path, "branch", "--list")
     if err != nil {
     	return nil, err
     }
     
     var localBranches []string
     localBranches = strings.Split(listBranches, "\n")

     for i, branch := range localBranches {
     	 branch = strings.TrimSpace(branch)
	 branch = strings.TrimPrefix(branch, "* ")
     	 localBranches[i] = branch
     }
     
     return localBranches, nil
}

func _RemoteBranches(apiResourcePrefix string) ([]string, error) {
     endpoint := apiResourcePrefix + "/branches"
     response, err := RunCommand("gh", "api", "-H", "Accept:application/vnd.github+json", endpoint)
     if err != nil {
     	return nil, err
     }

     var branches []map[string]string
     json.Unmarshal([]byte(response), &branches)

     var branchesNames = make([]string, len(branches))
     for i, branch := range branches {
     	 branchesNames[i] = branch["name"]
     }
     sort.Strings(branchesNames)

     return branchesNames, nil
}

func _Branches(localBranches []string, remoteBranches []string) []string {
     branchesMap := map[string]struct{}{}
     
     for _, branch := range localBranches {
     	 branchesMap[branch] = struct{}{}
     }
     for _, branch := range remoteBranches {
     	 branchesMap[branch] = struct{}{}
     }

     branches := make([]string, 0, len(branchesMap))
     for branch := range branchesMap {
     	 branches = append(branches, branch)
     }

     sort.Strings(branches)

     return branches
}

func _CurrentBranch(path string) (string, error) {
     return RunCommand("git", "-C", path, "branch", "--show-current")
}

func (rc RepositoryClone) HasRemoteBranch(branch string) (bool, error) {
     _, found := sort.Find(len(rc.RemoteBranches), func(i int) int {
     	return strings.Compare(branch, rc.RemoteBranches[i])
     })

     return found, nil
}

func (rc RepositoryClone) OpenPullRequest(jiraIssueId string,
					  pullRequestTemplate PullRequestTemplate,
					  sourceBranch string,
					  targetBranch string) (string, error) {
     endpoint := fmt.Sprintf("/repos/%s/%s/pulls", rc.RepoOwner, rc.Name)
	      
     response, err := RunCommand("gh", "api",
     	       "--method", "POST", "-H", "Accept:application/vnd.github+json", endpoint,
	       "-f", fmt.Sprintf("title=%s", jiraIssueId),
	       "-f", fmt.Sprintf("body=%s", pullRequestTemplate.Generate()),
	       "-f", fmt.Sprintf("head=%s", sourceBranch),
	       "-f", fmt.Sprintf("base=%s", targetBranch))
     if err != nil {
     	type PullRequestError struct {
 	     Message string
	}
	type PullRequestErrorResponse struct {
     	     Errors []PullRequestError
	}
	var pullRequestErrorResponse PullRequestErrorResponse
     	json.Unmarshal([]byte(err.Error()), &pullRequestErrorResponse)

	return "", errors.New(pullRequestErrorResponse.Errors[0].Message)
     }

     return UnmarshallJson(response)["html_url"], nil
}

func UnmarshallJson(value string) map[string]string {
     var pr map[string]string
     json.Unmarshal([]byte(value), &pr)
     return pr
}