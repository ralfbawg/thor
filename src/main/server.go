package main

import (
	"common/logging"
	"net/http"
)

func StartServers()  {
	logging.Info("start servers")

	
}
func startWsServer()  {
	http.HandleFunc("/", wsHandler)
	http.ListenAndServe("0.0.0.0:7777", nil)
}
func startApiServer()  {
	
}