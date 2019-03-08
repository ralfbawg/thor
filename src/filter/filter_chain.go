package filter

import (
	"common/logging"
	"net/http"
)

const (
	filterMax = 100
)

type filterI interface {
	before()
	do(w http.ResponseWriter, r *http.Request)
	after()
}

var filters = make([]filterI, 0, filterMax)

func FilterInit() {
	filters = append(filters, &AuthFilter{}) //TODO 未来需要换成配置方式加载
}

func DoFilter(w http.ResponseWriter, r *http.Request) {
	for _, t := range filters { //TODO 根据服务与任务配置不同过滤器
		logging.Debug("len=%d,cap=%d", len(filters), cap(filters))
		t.before()
		t.do(w, r)
		t.after()
	}
}
