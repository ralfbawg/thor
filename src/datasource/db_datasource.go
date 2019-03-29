package datasource

import (
	"database/sql"
	"common/logging"
	"github.com/go-xorm/xorm"
	"errors"
	"util"
)

const (
	DBTYPE_MYSQL  = "mysql"
	DBTYPE_ORACLE = "mysql"
)

type DbDatasource struct {
	DbName            string
	Username          string
	Password          string
	Url               string
	DbType            string
	Scheme            string
	MaxConnection     int
	MaxIdleConnection int
	Db                *sql.DB
	BaseDataSource
}

func (db *DbDatasource) init() {
	if datasource, err1 := db.GetDatasourceName(); err1 != nil {
		logging.Debug("err %s", err1.Error())
	} else if engine, err2 := xorm.NewEngine(db.DbType, datasource); err2 != nil {
		logging.Error("err %s", err2.Error())
	} else {
		dbSrc := engine.DB()
		dbSrc.SetMaxOpenConns(util.AOrB(func() bool { return db.MaxConnection <= 0 }, 2000, db.MaxConnection).(int))
		dbSrc.SetMaxIdleConns(util.AOrB(func() bool { return db.MaxIdleConnection <= 0 }, 1000, db.MaxIdleConnection).(int))
		dbSrc.Ping()
		db.Db = dbSrc.DB
		db.initFlag = true
	}

}

func (db *DbDatasource) GetConnection() *sql.DB {
	if !db.initFlag {
		db.init()
	}
	return db.Db
}

func getExtInfoByType(dbtype string) string {
	switch dbtype {
	case DBTYPE_MYSQL:
		return "?charset=utf8"
	default:
		return ""
	}
}

func (db *DbDatasource) GetDatasourceName() (string, error) {
	switch db.DbType {
	case "mysql":
		return db.Username + ":" + db.Password + "@tcp(" + db.Url + ")/" + db.DbName + getExtInfoByType(db.DbType), nil
	default:
		return "", errors.New("Get datasource error")

	}
}
