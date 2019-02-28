package api

import (
	"net/http"
	"strings"
	"fmt"
	"reflect"
)

type apiDispatch struct {
}

func dispatch() {

}
func (h *apiDispatch) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, world!"))
}

func startApiServer() {
	http.HandleFunc("/api/", &apiDispatch{})
	http.ListenAndServe(":80", nil)
}
func say(w http.ResponseWriter, req *http.Request) {
	pathInfo := strings.Trim(req.URL.Path, "/")
	fmt.Println("pathInfo:", pathInfo)

	parts := strings.Split(pathInfo, "/")
	fmt.Println("parts:", parts)

	var action = "ResAction"
	fmt.Println(strings.Join(parts, "|"))
	if len(parts) > 1 {
		fmt.Println("22222222")
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
