package api

import "net/http"

type apiDispatch struct {

}

func dispatch()  {
	
}
func (h *apiDispatch) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, world!"))
}
func startApiServer()  {
	 http.Handle("/api/",&apiDispatch{})
	 http.ListenAndServe("80",nil)
}
