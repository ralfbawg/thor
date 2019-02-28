package handler

import "net/http"

type WsHandler struct {
	uuid string
	httpHandler http.Handler

}

