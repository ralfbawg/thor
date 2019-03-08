package datasource

import "common/logging"

type DataSourceI interface {
	getData() interface{}
	init() interface{}
}
type BaseDataSource struct {
	inited bool
	DataSourceI
}

func (d *BaseDataSource) init() {
	logging.Debug("default init")
}

func (d *BaseDataSource) getData() {
	logging.Debug("test")
}
