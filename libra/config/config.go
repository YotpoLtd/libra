package config

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/hashicorp/hcl"
	log "github.com/sirupsen/logrus"
	"github.com/YotpoLtd/libra/libra/structs"
)

type Config struct {
	Jobs     map[string]*structs.NomadJob      `hcl:"job"`
	Nomad    structs.NomadConfig               `hcl:"nomad"`
	Backends map[string]structs.Backend `hcl:"backend"`
	Store    map[string]structs.Store   `hcl:"store"`
}

// NewConfig will return a Config struct
func ParseConfig(configDir string) (*Config, error) {

	info, err := os.Stat(configDir)
	if os.IsNotExist(err) {
		log.Errorf("Configuration Directory %s does not exists", configDir)
		return nil, err
	}

	if !info.IsDir() {
		log.Errorf("Specified Configuration Directory %s is not a directory", configDir)
	}

	fileList := []string{}
	err = filepath.Walk(configDir, func(filePath string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			fileList = append(fileList, filePath)
		}
		return nil
	})

	if len(fileList) == 0 {
		log.Errorf("Failed to find files in %s", configDir)
	}

	if err != nil {
		log.Errorf("Failed to detect file: %s", err)
		return nil, err
	}

	var configBlob bytes.Buffer
	for i, file := range fileList {
		log.Infof("File #%d: %s", i, file)
	}
	for _, file := range fileList {
		config, err := ioutil.ReadFile(file)
		if err != nil {
			log.Errorf("Failed to read file (%s): %s", file, err)
			return nil, err
		}
		configBlob.Write(config)
	}

	var out Config
	err = hcl.Decode(&out, configBlob.String())
	if err != nil {
		log.Errorf("HCL Error: %s", err)
		return nil, err
	}

	for jobName, jobConfig := range out.Jobs {
		jobConfig.Name = jobName

		for groupName, groupConfig := range jobConfig.Groups {
			groupConfig.Name = groupName

			for ruleName, ruleConfig := range groupConfig.Rules {
				ruleConfig.Name = ruleName
			}
		}
	}

	return &out, nil
}
