package config

import (
	"gopkg.in/yaml.v2"
	"log"
)

var ConfigStore = &Configure{}

type Configure struct {
	Common struct {
		Env string `yaml:"env"`
	}
	Log struct {
		Filepath string `yaml:"filepath"`
		Level string `yaml:"level"`
	}
	Server struct {
		Port           string `yaml:"port"`
		WsServicePath  string `yaml:"wspath"`
		ApiServicePath string `yaml:"apipath"`
	}
	Db struct {
		Host      string `yaml:"host"`
		Port   string `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		DbType   string `yaml:"dbtype"`
		DbName string `yaml:"dbname"`
	}
	inited bool
}

func (c *Configure) GetConfig(reload bool) (*Configure, error) {
	if !reload && ConfigStore.inited {
		return ConfigStore, nil
	}
	//yamlFile, err := ioutil.ReadFile("src/config.yaml")
	yamlFile, err := Asset("src/config/file/config.yaml")
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
