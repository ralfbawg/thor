package db

import (
	"common/logging"
	"datasource"
	"db/pojo"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var db *datasource.DbDatasource

func InitDb(host string, port string, dbname string, username string, password string, dbtype string) {
	db := &datasource.DbDatasource{
		Host:     host,
		Port:     port,
		DbName:   dbname,
		Username: username,
		Password: password,
		DbType:   dbtype,
	}
	db.Init()
	//db.GetDb().Sync2(new(pojo.TblThorApp))
	//result := &pojo.TblThorApp{}
	//db.GetDb().ShowSQL(true)

	result := make([]pojo.TblThorApp, 0)
	err := db.GetDb().Find(&result)
	if err != nil {
		fmt.Println(err)
	}
	logging.Debug("good")
}

func GetAppDb() *datasource.DbDatasource {
	if db == nil {
		//InitDb()
	}
	return db
}
