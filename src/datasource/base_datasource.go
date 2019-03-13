package datasource

import (
	"common/logging"
	"sync"
)

type DataSourceI interface {
	getData(id string) interface{}
	init() interface{}
}
type BaseDataSource struct {
	initFlag bool
	pool sync.Pool
	poolFlag bool
	DataSourceI
}

func (d *BaseDataSource) init() {
	d.initFlag = true
	logging.Debug("default init")
}

func (d *BaseDataSource) getData(id string) interface{} {
	logging.Debug("test")
	return "test"
}
