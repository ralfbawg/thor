package datasource

import "sync"

type MysqlDatasouce struct {
	DbDatasouce
}


func (mysql *MysqlDatasouce) init() {
	if mysql.usePool {
		mysql.pool = sync.Pool{
			New: mysql.getConnection,
		}
	}
}

func (mysql *MysqlDatasouce) getConnection() interface{} {

}
