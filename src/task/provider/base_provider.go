package provider

import "datasource"

type providerI interface {
	GetData(key interface{}) (interface{}, error)
}
type BaseProvider struct {
	datasource.BaseDataSource
	providerI
}
