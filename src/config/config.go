package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type configure struct {
	common struct{
		env string
	}
	log struct{
		filepath string
	}
	wserver struct{
		port int
	}
	apiserver struct{
		port int
	}
	db     struct {
		url      string `json"url"`
		username string `json"username"`
		password string `json"password"`
		dbtype   string `json"dbtype"`
	} `json:"db"`
}

func (c *configure) getConfig() (*configure, error) {
	yamlFile, err := ioutil.ReadFile("conf.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c, err
}
