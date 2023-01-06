package cmd

import (
       "fmt"
       "log"
       "os"
       "strings"
)

type ConfigFile struct {
     jiraServer	string
}

func InitConfigFile(jiraServer string) ConfigFile {
     return ConfigFile {
     	    jiraServer: jiraServer,
     }
}

func ReadConfigFile() (*ConfigFile, error) {
     data, err := os.ReadFile(ConfigFilePath())
     if err != nil {
     	return nil, err
     }
     
     return &ConfigFile {
     	    jiraServer: strings.Split(string(data), "=")[1],
     }, nil
}

func (this ConfigFile) Persist() {
     configFilePath := ConfigFilePath()
     _, err := os.Stat(configFilePath)
     if err != nil {
     	configDirectory := ConfigDirectory()
     	err = os.Mkdir(configDirectory, os.ModePerm)
	if (err != nil) {
     	   log.Fatal(err.Error())
     	}

     	_, err:= os.Create(configFilePath)
     	if (err != nil) {
     	   log.Fatal(err.Error())
     	}
     }

     data := []byte(fmt.Sprintf("jiraServer=%s", this.jiraServer))
     err = os.WriteFile(configFilePath, data, os.ModePerm)
     if err != nil {
     	log.Fatal(err.Error())
     }
}

func ConfigDirectory() string {
     homeDirectory, err := os.UserHomeDir()
     if err != nil {
     	log.Fatal(err)
     }
     return fmt.Sprintf("%s/.config/.gcpout", homeDirectory)
}

func ConfigFilePath() string {
     return fmt.Sprintf("%s/.config.properties", ConfigDirectory())
}