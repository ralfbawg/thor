package datasource

import (
	"common/logging"
	"errors"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
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
	Host              string
	Port              string
	DbType            string
	Scheme            string
	MaxConnection     int
	MaxIdleConnection int
	Db                *xorm.Engine
	BaseDataSource
}

func (db *DbDatasource) Init() (*DbDatasource, error) {
	if datasource, err1 := db.GetDatasourceName(); err1 != nil {
		logging.Debug("err %s", err1.Error())
		return nil, err1
	} else if engine, err2 := xorm.NewEngine(db.DbType, datasource); err2 != nil {
		logging.Error("err %s", err2.Error())
	} else {
		engine.SetTableMapper(core.SnakeMapper{}) //table名称为驼峰
		engine.SetColumnMapper(core.SameMapper{}) //列名对应为名字相同
		dbSrc := engine.DB()
		dbSrc.SetMaxOpenConns(util.AOrB(func() bool { return db.MaxConnection <= 0 }, 2000, db.MaxConnection).(int))
		dbSrc.SetMaxIdleConns(util.AOrB(func() bool { return db.MaxIdleConnection <= 0 }, 1000, db.MaxIdleConnection).(int))
		dbSrc.Ping()
		engine.ShowSQL(true)
		db.Db = engine
		db.initFlag = true
	}
	return db, nil
}

func (db *DbDatasource) GetDb() *xorm.Engine {
	if !db.initFlag {
		db.Init()
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

func (db *DbDatasource) getExtInfoByType() string {
	switch db.DbType {
	case DBTYPE_MYSQL:
		return "?charset=utf8"
	default:
		return ""
	}
}

func (db *DbDatasource) GetDatasourceName() (string, error) {
	switch db.DbType {
	case DBTYPE_MYSQL:
		return db.Username + ":" + db.Password + "@tcp(" + db.Host + ":" + db.Port + ")/" + db.DbName + db.getExtInfoByType(), nil
		//return db.Username + ":" + db.Password + "@tcp(" + db.Host + ":" + db.Port + ")/" + db.DbName + getExtInfoByType(db.DbType), nil
	default:
		return "", errors.New("don't support" + db.DbType + " %s yet")

	}
}

func (d *DbDatasource) getData(param ...string) interface{} {
	logging.Debug("test")
	return "test"
}
