package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

var ConfigStore = &Configure{}

type Configure struct {
	Common struct {
		Env string `yaml:"env"`
	}
	Log struct {
		Filepath string `yaml:"filepath"`
	}
	Server struct {
		Port           string `yaml:"port"`
		WsServicePath  string `yaml:"wspath"`
		ApiServicePath string `yaml:"apipath"`
	}
	Db struct {
		Url      string `yaml:"url"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Dbtype   string `yaml:"dbtype"`
	}
	inited bool
}

func (c *Configure) GetConfig(reload bool) (*Configure, error) {
	if !reload && ConfigStore.inited {
		return ConfigStore, nil
	}
	yamlFile, err := ioutil.ReadFile("src/config.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	c.inited = true
	return c, err
}
