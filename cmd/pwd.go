package cmd

import (
       "strings"
)

func Pwd() (string, error) {
     return ReadLink(".")
}

func ReadLink(path string) (string, error) {
     return RunCommand("readlink", "-f", path)
}

func LastPath(path string) string {
     hierarchy := strings.Split(path, "/")
     return hierarchy[len(hierarchy)-1]
}