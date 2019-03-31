package datasource

import (
	"common/logging"
	"sync"
)

type DataSourceI interface {
	GGetData(param ...string) interface{}
	Init() interface{}
}
type BaseDataSource struct {
	initFlag bool
	pool     sync.Pool
	poolFlag bool
	DataSourceI
}

func (d *BaseDataSource) Init() {
	d.initFlag = true
	logging.Debug("default init")
}

func (d *BaseDataSource) GetData(param ...string) interface{} {
	logging.Debug("test")
	return "test"
}
