package filter

import (
	"common/logging"
	"net/http"
)

type apiFilterI interface {
	before()
	after()
	do(w http.ResponseWriter, r *http.Request) bool
}
type wsFilterI interface {
	before(msg []byte)
	do(msg []byte) bool
	after(msg []byte)
}
type BaseApiFilter struct {
	apiFilterI
}

type BaseWsFilter struct {
	wsFilterI
}

func (b *BaseApiFilter) before() { //空方法
	logging.Debug("do before")
}
func (b *BaseApiFilter) after() { //空方法
	logging.Debug("do after")
}

func RegisterApiFilter(name string, f *BaseApiFilter) {
	if f == nil {
		panic("core: Register filter is nil")
	}
	if _, dup := ApiFilters[name]; dup {
		panic("core: Register called twice for filter " + name)
	}
	ApiFilters[name] = f

}
func RegisterWsFilter(name string, f *BaseWsFilter) { //空方法
	if f == nil {
		panic("core: Register filter is nil")
	}
	if _, dup := WsFilters[name]; dup {
		panic("core: Register called twice for filter " + name)
	}
	WsFilters[name] = f

}
func GetWsFilter(appId string) *BaseWsFilter {
	return WsFilters[appId].(*BaseWsFilter )
}
func WsFilterSize() int {
	return len(WsFilters)
}
