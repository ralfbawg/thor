package datasource

type RedisDatasource struct {
	BaseDataSource
}

func (r *RedisDatasource)getData()interface{}  {
	return "Test"
}
