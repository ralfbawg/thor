package filter

import "net/http"

const (
	filterMax = 100
)

var (
	ApiFilters = make(map[string]*BaseApiFilter, filterMax)
	WsFilters  = make(map[string]*BaseWsFilter, filterMax)
)

func DoApiFilter(w *http.ResponseWriter, r *http.Request) {
	for k, v := range ApiFilters {
		if k != "" {
			v.before()
			v.do(w, r)
			v.after()
		}
	}
}
func DoWsFilter(msg []byte) {
	for k, v := range WsFilters {
		if k != "" {
			v.before(msg)
			v.do(msg)
			v.after(msg)
		}
	}
}

type WsFilterChanin []*BaseWsFilter
type ApiFilterChanin []*BaseApiFilter

func NewFilterChain(appId string) WsFilterChanin {

	return nil
}
func findFilterByAppId() {

}
func getDefaultFilter() {

}

func (fc WsFilterChanin) doWsFilter(msg []byte) {
	for _, v := range fc {
		v.before(msg)
		v.do(msg)
		v.after(msg)
	}
}
