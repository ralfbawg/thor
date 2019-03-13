package datasource

import "database/sql"

type DbDatasouce struct {
	BaseDataSource
	username string
	password string
	url      string
	dbType   string
	scheme   string
	maxConnection int
	maxIdleConnection int
}
func (db *DbDatasouce)init() {
	dbSrc, _ = sql.Open(db.dbType, "root:@tcp(127.0.0.1:3306)/test?charset=utf8")
	dbSrc.SetMaxOpenConns(2000)
	dbSrc.SetMaxIdleConns(1000)
	dbSrc.Ping()
}
func GetDbSource(dbType string)*DbDatasouce  {
	dbSrc, _ = sql.Open(dbType, "root:@tcp(127.0.0.1:3306)/test?charset=utf8")
	dbSrc.SetMaxOpenConns(2000)
	dbSrc.SetMaxIdleConns(1000)
	dbSrc.Ping()
}

func (db *DbDatasouce)GetConnect() {
	dbSrc, _ = sql.Open(db.dbType, "root:@tcp(127.0.0.1:3306)/test?charset=utf8")
	dbSrc.SetMaxOpenConns(2000)
	dbSrc.SetMaxIdleConns(1000)
	dbSrc.Ping()
}
