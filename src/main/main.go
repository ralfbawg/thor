package main

import (
	"common/logging"
	"config"
)



func main() {
	config.Init_main()
	logging.Debug("server start")

}
