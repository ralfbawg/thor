package filter

import "net/http"

const (
	filterMax            = 100
	appFilterDefaultSize = 3
)

var (
	ApiFilters = make(map[string]interface{}, filterMax)
	WsFilters  = make(map[string]interface{}, filterMax)
)

func DoApiFilter(w http.ResponseWriter, r *http.Request) {
	for k, v := range ApiFilters {
		if k != "" {
			t := v.(BaseApiFilter)
			t.before()
			t.do(w, r)
			t.after()
		}
	}
}
func DoWsFilter(msg []byte) {
	for k, v := range WsFilters {
		if k != "" {
			t := v.(BaseWsFilter)
			t.before(msg)
			t.do(msg)
			t.after(msg)
		}
	}
}

type WsFilterChain []interface{}
type ApiFilterChain []interface{}

func NewWsFilterChain(appId string) WsFilterChain {
	t := make(WsFilterChain, appFilterDefaultSize)
	reg := &RegFilter{}
	t = append(t, reg)
	return t
}
func findFilterByAppId() {

}
func getDefaultFilter() {

}

func (fc WsFilterChain) doWsFilter(msg []byte) {
	for _, t := range fc {
		v := t.(wsFilterI)
		v.before(msg)
		v.do(msg)
		v.after(msg)
	}
}
