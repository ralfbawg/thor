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

func (nosql *NosqlDatasource) Init() {

}

func (nosql *NosqlDatasource) GetData(param interface{}) interface{} {
	return nil
}
