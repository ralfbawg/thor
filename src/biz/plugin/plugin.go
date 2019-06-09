package plugin

import (
	"fmt"
	"os"
	"plugin"
	"util"
)

type pluginI interface {
	init(interface{}) (bool, error)
	do(interface{}) interface{}
}
type Plugin struct {
	name string
	path string
}

type PluginManager struct {
	plugins  util.ConcMap
	register chan *plugin.Plugin
}

func (pm PluginManager) Register(name string, path string) {
	defer func() {
		recover()
	}()
	p, err := plugin.Open(path)
	if err != nil {
		fmt.Println("error open plugin: ", err)
		os.Exit(-1)
	}
	p.Lookup("init")
}
