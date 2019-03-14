package datasource

import "sync"

type MysqlDatasource struct {
	DbDatasource
}

func (mysql *MysqlDatasource) init() {
	if mysql.poolFlag {
		mysql.pool = sync.Pool{
			New: mysql.getConnection,
		}
	}
}

func (mysql *MysqlDatasource) getConnection() interface{} {
	return nil
}
