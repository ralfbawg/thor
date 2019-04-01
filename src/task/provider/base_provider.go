package provider

import (
	"github.com/panjf2000/ants"
	"datasource"
)

var(
	providerPoolDefaultSize = 100
	providerPool,_ = ants.NewPool(providerPoolDefaultSize)
)

type providerI interface {
	Init()
	Run()
	GetData(key interface{}) (interface{}, error)
}
type BaseProvider struct {
	ds *datasource.BaseDataSource
	providerI
}
func (p *BaseProvider) Init() {
	p.ds.Init()
}

func (p *BaseProvider) Run() {
	providerPool.Submit(func() {

	})
}
