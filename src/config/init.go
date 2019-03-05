package config

import (
	"common/logging"
	"db"
)
func Init_main() {
	initConfigFile()
	db.InitDb()
}
func initConfigFile() (*Configure,error) {
	logging.Debug("init db start")
	if c,err := ConfigStore.GetConfig(true);err==nil{
		logging.Debug("db init success")
		return c,nil
	}else {
		logging.Debug("db init fail")
		return nil,err
	}

}
func GetConfigFile() (*Configure,error) {
	return ConfigStore,nil
}

