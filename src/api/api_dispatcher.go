package api

import (
	"net/http"
	"strings"
	"io"
	"io/ioutil"
	ws "websocket"

)

const SuccessMsg = "success"

func ApiDispatch(w http.ResponseWriter, r *http.Request) {
	paths := strings.Split(r.RequestURI, "/")
	if strings.HasPrefix(paths[2], "b") {
		msg := r.URL.Query().Get("msg")
		appId := r.URL.Query().Get("appId")
		uid := r.URL.Query().Get("uid")
		go ws.WsBroadcast(appId, uid, msg)
		w.Write([]byte(SuccessMsg))
	}
	io.Copy(ioutil.Discard, r.Body)
	r.Body.Close()
}
