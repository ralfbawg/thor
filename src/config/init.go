package config

import "common/logging"

func Init_main() {
	init_configFile()
}
func init_configFile() {
	logging.Debug("init main")
}

