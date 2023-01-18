package cmd

import (
       "path/filepath"
       "os"
       "strings"
)

type Dir struct {
     pwd	string
     children	[]string
}

func ProjectsDir() (*Dir, error) {
     var inner func(relativePath string) (string, error)
     inner = func (relativePath string) (string, error) {
     	     	  if IsGitRepository(relativePath) {
     	       	     return inner(relativePath + "../")
     	     	  }
     	     	  return RunCommand("readlink", "-f", relativePath)
     	     }
     path, err := inner("")
     if err != nil {
     	return nil, err
     }
     return Ls(path + "/"), nil
}

func Ls(path string) *Dir {
     dir := &Dir{}
     dir.pwd = path
     var pathResults, _ = filepath.Glob(path + "*")
     dir.children = make([]string, len(pathResults))
     for i, v := range pathResults {
     	 hierarchy := strings.Split(v, "/")
	 dir.children[i] = hierarchy[len(hierarchy)-1]
     }

     return dir
}

func (this Dir) GitRepositories() *Dir {
     var gitRepositories []string = make([]string, len(this.children))
     i := 0
     for _, child := range this.children {
     	 if IsGitRepository(this.pwd + child + "/") {
	    gitRepositories[i] = child
	    i += 1
	 }
     }
     gitRepositories = gitRepositories[:i]

     return &Dir {
     	    pwd: this.pwd,
	    children: gitRepositories,
     }
}

func IsGitRepository(path string) bool {
     _, err := os.Stat(path + ".git")
     return err == nil
}