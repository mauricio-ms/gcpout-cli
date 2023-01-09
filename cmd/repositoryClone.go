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
}

// TODO Create New method to precompute the methods into fields

func (rc RepositoryClone) Path() string {
     return rc.ParentPath + "/" + rc.Name
}

func (rc RepositoryClone) CurrentBranch() (string, error) {
     return RunCommand("git", "-C", rc.Path(), "branch", "--show-current")
}

func (rc RepositoryClone) RepoOwner() (string, error) {
     originUrl, err := RunCommand("git", "-C", rc.Path(), "remote", "get-url", "origin")
     if err != nil {
     	return "", err
     }
     
     ownerRegex := regexp.MustCompile(`^[^:]+:([^/]+).*$`)
     return ownerRegex.FindStringSubmatch(originUrl)[1], nil
}

func (rc RepositoryClone) Branches() ([]string, error) {
     branchesMap := map[string]struct{}{}
     
     localBranches, err := rc.LocalBranches()
     if err != nil {
     	return nil, err
     }
     for _, branch := range localBranches {
     	 branchesMap[branch] = struct{}{}
     }

     remoteBranches, err := rc.RemoteBranches()
     if err != nil {
     	return nil, err
     }
     for _, branch := range remoteBranches {
     	 branchesMap[branch] = struct{}{}
     }

     branches := make([]string, 0, len(branchesMap))
     for branch := range branchesMap {
     	 branches = append(branches, branch)
     }

     sort.Strings(branches)

     return branches, nil
}

func (rc RepositoryClone) LocalBranches() ([]string, error) {
     listBranches, err := RunCommand("git", "-C", rc.Path(), "branch", "--list")
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

func (rc RepositoryClone) HasRemoteBranch(branch string) (bool, error) {
     remoteBranches, err := rc.RemoteBranches()
     if err != nil {
     	return false, err
     }
     sort.Strings(remoteBranches)

     _, found := sort.Find(len(remoteBranches), func(i int) int {
     	return strings.Compare(branch, remoteBranches[i])
     })

     return found, nil
}

func (rc RepositoryClone) RemoteBranches() ([]string, error) {
     apiResourcePrefix, err := rc.ApiResourcePrefix()
     if err != nil {
     	return nil, err
     }
     
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

     return branchesNames, nil
}

func (rc RepositoryClone) ApiResourcePrefix() (string, error) {
     repoOwner, err := rc.RepoOwner()
     if err != nil {
     	return "", nil
     }
     
     return fmt.Sprintf("/repos/%s/%s", repoOwner, rc.Name), nil
}

func (rc RepositoryClone) OpenPullRequest(jiraIssueId string,
					  pullRequestTemplate PullRequestTemplate,
					  sourceBranch string,
					  targetBranch string) (string, error) {
     repoOwner, err := rc.RepoOwner()
     if err != nil {
     	return "", nil
     }
     endpoint := fmt.Sprintf("/repos/%s/%s/pulls", repoOwner, rc.Name)
	      
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