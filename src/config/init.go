package config

import (
	"common/logging"
	"fmt"
)

func Init_main() {
	init_configFile()
}
func init_configFile() {
	logging.Debug("init main")
	configure := &configure{}
	if c,err := configure.getConfig();err!=nil{
		return
	}else {
		fmt.Printf(c.Db.Dbtype)
	}


}

