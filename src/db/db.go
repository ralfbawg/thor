package db

import (
	"common/logging"
	"datasource"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

var db *datasource.DbDatasource

func InitDb() {
	db := &datasource.DbDatasource{
		DbType:            "mysql",
		Url:               "14.17.108.89:4001",
		DbName:            "yyweb_new",
		Username:          "yymz",
		Password:          "duowan",
		MaxConnection:     100,
		MaxIdleConnection: 100,
	}
	engine, err := xorm.NewEngine(db.DbType, get)
	//if rows, err := db.GetConnection().Query("select * from TBL_CI_REPO"); err == nil {
	//	columns, _ := rows.Columns()
	//	scanArgs := make([]interface{}, len(columns))
	//	values := make([]interface{}, len(columns))
	//	for i := range values {
	//		scanArgs[i] = &values[i]
	//	}
	//	for rows.Next() {
	//		//将行数据保存到record字典
	//		err = rows.Scan(scanArgs...)
	//		record := make(map[string]string)
	//		for i, col := range values {
	//			if col != nil {
	//				record[columns[i]] = string(col.([]byte))
	//			}
	//		}
	//		fmt.Println(record)
	//	}
	//} else {
	//	logging.Debug("error occured %s", err)
	//}
	logging.Debug("good")
}
