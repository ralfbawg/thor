package api

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

type apiServer struct {
	mainHandler *http.Handler
}

func (c *apiServer) ApiService(w *http.ResponseWriter,r *http.Request)  {
	pathInfo := strings.Trim(r.URL.Path, "/")
	fmt.Println("pathInfo:", pathInfo)

	parts := strings.Split(pathInfo, "/")
	fmt.Println("parts:", parts)

	var action = "ResAction"
	fmt.Println(strings.Join(parts, "|"))
	if len(parts) > 1 {
		action = strings.Title(parts[1]) + "Action"
	}
	fmt.Println("action:", action)
	handle := &Handlers{}
	controller := reflect.ValueOf(handle)
	println(controller.Bytes())
	method := controller.MethodByName(action)
	r := reflect.ValueOf(req)
	wr := reflect.ValueOf(w)
	method.Call([]reflect.Value{wr, r})
}
