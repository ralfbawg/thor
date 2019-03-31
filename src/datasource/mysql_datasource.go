package datasource

type MysqlDatasource struct {
	*DbDatasource
}

func NewMysql(ip string, port string) *DbDatasource {
	return &DbDatasource{
		DbType: DBTYPE_MYSQL,
		Host:   ip,
		Port:   port,
	}
}
func (mysql MysqlDatasource) getExtInfoByType() string {
	return "?charset=utf8"
}
