package datasource

type RedisDatasource struct {
	BaseDataSource
}

func (r *RedisDatasource) getData(id string) interface{} {
	return "Test"
}
func (r *RedisDatasource) init() {

}
