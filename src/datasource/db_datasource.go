package datasource

import (
	"database/sql"
	"common/logging"
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
	//if engine, err := xorm.NewEngine(db.DbType, db.Username+":"+db.Password+"@tcp("+db.Url+")/"+db.DbName+getExtInfoByType(db.DbType)); err != nil {
	if dbSrc, err := sql.Open(db.DbType, db.Username+":"+db.Password+"@tcp("+db.Url+")/"+db.DbName+getExtInfoByType(db.DbType)); err != nil {
		logging.Debug("err %s", err)
	} else {
		//engine.DB()
		dbSrc.SetMaxOpenConns(2000)
		dbSrc.SetMaxIdleConns(1000)
		dbSrc.Ping()
		db.Db = dbSrc
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
	case "mysql":
		return "?charset=utf8"
	default:
		return ""
	}
}

func (db *DbDatasource) GetDatasourceName(dbtype string) (string, error) {
	switch dbtype {
	case "mysql":
		return db.Username + ":" + db.Password + "@tcp(" + db.Url + ")/" + db.DbName + getExtInfoByType(db.DbType), nil
	default:
		return "", new(error)

	}
}
