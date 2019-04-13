package config

import (
	"common/logging"
	"filter"
	"monitor"
	"websocket"
)

func InitMain() {
	if c, error := initConfigFile(); error != nil {

	} else {
		logging.Debug("db hots:%s", c.Db.Host)
		//db.InitDb(c.Db.Host, c.Db.Port, c.Db.DbName, c.Db.Username, c.Db.Password, c.Db.DbType)
		monitor.MonitorInit()
		filter.FilterInit()
		websocket.WsManagerInit()
	}

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
