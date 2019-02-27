package main

import (
	"common/logging"
	"net/http"
)

const ()

func StartServers() {
	logging.Info("start servers")
	startWsServer()

}
func startWsServer() {
	http.HandleFunc("/", wsHandler)
	http.ListenAndServe("0.0.0.0:7777", nil)
}
func startApiServer() {

}
