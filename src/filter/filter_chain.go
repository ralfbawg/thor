package filter

import (
	"net/http"
	"fmt"
)

type filterI interface {
	before()
	do(w http.ResponseWriter,r *http.Request)
	after()
}

var filters = make([]filterI,0,10)

func FilterInit()  {
	filters = append(filters,&AuthFilter{&BaseFilter{},})//TODO 未来需要换成配置方式加载
}


type FilterChain struct {
	w http.ResponseWriter
	r *http.Request
	filters []filterI
}

func DoFilter(w http.ResponseWriter,r *http.Request) {
	for _,t := range filters{
		fmt.Println("%d,%d",len(filters),cap(filters))
		fmt.Printf("%p, %T\n", t, t)
		t.before()
		t.do(w,r)
		t.after()
	}
}
