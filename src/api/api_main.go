package api

import (
	"net/http"
	"util/uuid"
	"strings"
	"crypto"
	"time"
	"errors"
	"strconv"
)

const (
	defaultSignAlgorithm = "SHA1"
)

type apiServer struct {
	mainHandler *http.Handler
}

type AppInfo struct {
	AppId           string
	AppKey          string
	AppSecret       string
	Desc            string
	cryptoAlgorithm string //加密算法
}

func (c *apiServer) ApiService(w *http.ResponseWriter, r *http.Request) {

}

//生成app信息，包括key与secret
func GernateAppInfo(appId string) (*AppInfo, error) {
	appKey := strings.ReplaceAll(uuid.Generate().String(), "-", "")
	h := crypto.SHA1.New()
	_, error := h.Write([]byte(appId + appKey))
	appSecret := h.Sum(nil)
	if error != nil {
		return nil, error
	}
	return &AppInfo{AppId: appId, AppKey: appKey, AppSecret: strings.ToUpper(string(appSecret)), cryptoAlgorithm: defaultSignAlgorithm}, nil
}

//验证App信息
func VerifyAppInfo(appId string, sign string, timestamp int64, random int, appSecret string) (error, bool) {
	if timestamp-time.Now().Unix() > int64(1*time.Minute) {
		return errors.New("timestamp is unllegal"), false
	} else if CheckSign(&AppInfo{AppId: appId, AppKey: appId, AppSecret: appSecret}, sign, random, timestamp) {
		return nil, true
	} else {
		return errors.New("sign is unllegal"), false
	}
}

func CheckSign(app *AppInfo, sign string, random int, timestamp int64) bool {
	_sign, _ := Sign(app.AppKey, app.AppSecret, random, timestamp, getHashAlgorithm(app.cryptoAlgorithm))
	if string(_sign) == sign {
		return true
	} else {
		return false
	}
}

func Sign(appkey string, appSecret string, random int, timestamp int64, hash crypto.Hash) ([]byte, error) {
	a := appkey + appSecret + strconv.Itoa(random) + strconv.Itoa(int(timestamp))
	h := hash.New()
	_, error := h.Write([]byte(a))
	sign := h.Sum(nil)
	if error == nil {
		return sign, nil
	} else {
		return nil, error
	}
}

func getHashAlgorithm(name string) (crypto.Hash) {
	switch strings.ToUpper(name) {
	case "SHA1":
		return crypto.SHA1
	case "MD5":
		return crypto.MD5
	case "MD5SHA1":
		return crypto.MD5SHA1
	case "SHA256":
		return crypto.SHA256
	default:
		return crypto.SHA1
	}
}
