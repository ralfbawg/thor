package api

import (
	"net/http"
)

type apiServer struct {
	mainHandler *http.Handler
}

func (c *apiServer) ApiService(w *http.ResponseWriter, r *http.Request) {

}
