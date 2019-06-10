package api

import (
	"net/http"
	"util/uuid"
	"strings"
	"crypto"
)

type apiServer struct {
	mainHandler *http.Handler
}

type AppInfo struct {
	AppId     string
	AppKey    string
	AppSecret string
	Desc      string
}

func (c *apiServer) ApiService(w *http.ResponseWriter, r *http.Request) {

}

func GernateAppInfo(appId string) (*AppInfo, error) {
	appKey := strings.ReplaceAll(uuid.Generate().String(), "-", "")
	h := crypto.SHA1.New()
	_, error := h.Write([]byte(appId + appKey))
	appSecret := h.Sum(nil)
	if error != nil {
		return nil, error
	}
	return &AppInfo{AppId: appId, AppKey: appKey, AppSecret: string(appSecret),}, nil
}

func VerifyAppInfo(appId string, sign string, timestamp int64, random int) {

}
