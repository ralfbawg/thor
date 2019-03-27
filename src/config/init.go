package config

import (
	"common/logging"
	"filter"
	"websocket"
	"monitor"
)

func InitMain() {
	initConfigFile()
	//db.InitDb()
	monitor.MonitorInit()
	filter.FilterInit()
	websocket.WsManagerInit()
}
func initConfigFile() (*Configure, error) {
	logging.Debug("init db start")
	if c, err := ConfigStore.GetConfig(true); err == nil {
		logging.Init(c.Log.Level)
		logging.Debug("db init success")
		return c, nil
	} else {
		logging.Debug("db init fail")
		return nil, err
	}

}
func GetConfigFile() (*Configure, error) {
	return ConfigStore, nil
}
