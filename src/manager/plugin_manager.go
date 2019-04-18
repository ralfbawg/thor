package manager

import "github.com/panjf2000/ants"

var PluginM = &PluginManager{}

type PluginI interface {
	Init(param interface{}) bool
}

type WsPluginI interface {
	Run()
	GetData(id interface{})
	PluginI
}
type BaseWsPlugin struct {
	WsPluginI
}
type PluginManager struct{}

func InitWsPlugin(plugins []*BaseWsPlugin) {
	for _, v := range plugins {
		ants.Submit(v.Run)
	}
}
func (pm *PluginManager) InitPlugin() {
	InitWsPlugin(nil)
}
