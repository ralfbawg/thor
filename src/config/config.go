package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type configure struct {
	Common struct{
		Env string `yaml:"env"`
	}
	Log struct{
		Filepath string `yaml:"filepath"`
	}
	server struct{
		Port int `yaml:"port"`
		WsServicePath string `yaml:"wspath"`
		ApiServicePath string  `yaml:"apipath"`
	}
	Db     struct {
		Url      string `yaml:"url"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		Dbtype   string `yaml:"dbtype"`
	}
	common struct{
		env string
	}`yaml:"common"`
	log struct{
		filepath string
	}`yaml:"log"`
	wserver struct{
		port int
	}`yaml:"wsserver"`
	apiserver struct{
		port int
	}`yaml:"apiserver"`
	db     struct {
		url      string `yaml"url"`
		username string `yaml"username"`
		password string `yaml"password"`
		dbtype   string `yaml"dbtype"`
	} `yaml:"db"`
}

func (c *configure) getConfig() (*configure, error) {
	yamlFile, err := ioutil.ReadFile("config.yaml")
	ioutil.ReadDir()
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
return c, err
}
