package config

import (
	"common/logging"
	"comet/websocket"
	"task"
)

func InitMain() {
	if c, error := initConfigFile(); error != nil {

	} else {
		logging.Debug("db hots:%s", c.Db.Host)
		//db.InitDb(c.Db.Host, c.Db.Port, c.Db.DbName, c.Db.Username, c.Db.Password, c.Db.DbType)
		//monitor.MonitorInit()
		websocket.WsManagerInit()
		task.AppManagerInst.SetAppInfo(&task.AppInfo{
			AppId:     ConfigStore.Ws.App.AppId,
			AppSecret: ConfigStore.Ws.App.AppSecret,
			AppKey:    ConfigStore.Ws.App.AppKey,
			Desc:      "",
		})
	}

}
func initConfigFile() (*Configure, error) {
	logging.Debug("init config start")
	if c, err := ConfigStore.GetConfig(true); err == nil {
		logging.Init(c.Log.Level, c.Log.Filepath)
		logging.Debug("config init success")
		return c, nil
	} else {
		logging.Debug("config init fail")
		return nil, err
	}

}
func GetConfigFile() (*Configure, error) {
	return ConfigStore, nil
}
