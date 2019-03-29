package datasource

type MysqlDatasource struct {
	DbDatasource
}

func NewMysql(ip string, port string) *DbDatasource {
	return &DbDatasource{
		DbType: DBTYPE_MYSQL,
		Url:    ip + ":" + port,
	}
}
