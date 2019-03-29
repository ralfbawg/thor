package datasource

type NosqlDatasource struct {
	DbName            string
	Username          string
	Password          string
	Url               string
	DbType            string
	Scheme            string
	MaxConnection     int
	MaxIdleConnection int
	BaseDataSource
}

