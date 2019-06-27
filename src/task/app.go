package task

import (
	"net/url"
	"util/uuid"
	"strconv"
	"common/logging"
	"strings"
	"crypto/sha1"
	"io"
	"fmt"
	"errors"
	"crypto"
	"util"
	"time"
)

const (
	appIdParam     = "appId"
	taskIdParam    = "taskId"
	appKeyParam    = "appKey"
	appSecretParam = "appSecret"
	randonParam    = "random"
	timestampParam = "timestamp"
	signParam      = "sign"
	uidParam       = "uid"
)

const (
	defaultSignAlgorithm    = "SHA1"
	defaultTimestampWinTime = int64(3 * time.Minute / 1e6)
)

var AppManagerInst = Init()
/*
 验证app信息
*/
func VerifyAppInfo(param url.Values) (string, int, string, bool, error) {
	appId := param.Get(appIdParam)
	appInfo, err := AppManagerInst.GetAppInfo(appId)
	if err != nil {
		return "", 0, "", false, err
	}
	//appSecret := param.Get(appSecretParam)
	random := param.Get(randonParam)
	timestamp := param.Get(timestampParam)
	sign := param.Get(signParam)
	uid := param.Get(uidParam)
	//logging.Debug("current param is appId:%s appKey:%s uid:%s", appId, appKey, uid)
	taskId, err := strconv.Atoi(param.Get(taskIdParam))
	if err != nil {
		taskId = 0
	}
	if uid == "" {
		uid = uuid.Generate().String()
	}
	logging.Debug("websocket connected,app(%s),appkey(%s),uid (%s)", appId, appInfo.AppKey, uid)
	return VerifyAppInfo2(appId, taskId, uid, sign, random, timestamp, appInfo.AppSecret)
}

/*
 验证app信息
*/
func VerifyAppInfo2(appId string, taskId int, uid string, sign string, random string, timestamp string, appSecret string) (string, int, string, bool, error) {
	//logging.Debug("app id is %s,app key is %s,uid is %s", appId, appKey, id)
	error, ok := VerifyAppInfo3(appId, sign, timestamp, random, appSecret)
	if ok && error == nil {
		return appId, taskId, uid, true, error
	} else {
		return appId, taskId, uid, false, error
	}
	return appId, taskId, uid, true, error
}

type AppManager struct {
	apps util.ConcMap
}

func Init() *AppManager {
	tmp := &AppManager{util.NewConcMap()}
	tmp.LoadAllAppinfos()
	return tmp
}

func (a *AppManager) LoadAllAppinfos() (bool, error) { //TODO 未来改为db式
	return true, nil
}
func (a *AppManager) SetAppInfo(info *AppInfo) (bool, error) { //TODO 未来改为db式
	a.apps.Set(info.AppId, info)
	return true, nil
}
func (a *AppManager) GetAppInfo(appId string) (*AppInfo, error) {
	tmp, exist := a.apps.Get(appId)
	if exist {
		return tmp.(*AppInfo), nil
	} else {
		return nil, errors.New("不存在")
	}
}

type AppInfo struct {
	AppId           string
	AppKey          string
	AppSecret       string
	Desc            string
	cryptoAlgorithm string //加密算法
}

//生成app信息，包括key与secret
func GernateAppInfo(appId string) (*AppInfo, error) {
	appKey := strings.ReplaceAll(uuid.Generate().String(), "-", "")
	h := sha1.New()
	encodeStr := appId + appKey
	io.WriteString(h, encodeStr)
	test := h.Sum(nil)
	appSecret := fmt.Sprintf("%x", test)
	return &AppInfo{AppId: appId, AppKey: appKey, AppSecret: appSecret, cryptoAlgorithm: defaultSignAlgorithm}, nil
	//return nil, nil
}

//验证App信息
func VerifyAppInfo3(appId string, sign string, timestamp string, random string, appSecret string) (error, bool) {
	//return nil, true //TODO 测试用
	timestampInt, covError1 := strconv.Atoi(timestamp)
	randomInt, covError2 := strconv.Atoi(random)
	if covError1 != nil || covError2 != nil {
		trueError := util.AOrB(func() bool { return covError1 != nil }, covError1, covError2)
		return errors.New("param conv is unllegal(" + trueError.(error).Error() + ")"), false
	}
	now := time.Now().Unix() * 1000
	if now-int64(timestampInt) > defaultTimestampWinTime || now-int64(timestampInt) < -defaultTimestampWinTime {
		return errors.New("timestamp is unllegal"), false
	} else if CheckSign(&AppInfo{AppId: appId, AppKey: appId, AppSecret: appSecret}, sign, randomInt, int64(timestampInt)) {
		return nil, true
	} else {
		return errors.New("sign is unllegal"), false
	}
}

func CheckSign(app *AppInfo, sign string, random int, timestamp int64) bool {
	_sign, _ := Sign(app.AppKey, app.AppSecret, random, timestamp, getHashAlgorithm(app.cryptoAlgorithm))
	logging.Debug("post sign=%s,calc sign=%s", sign, _sign)
	if _sign == sign {
		return true
	} else {
		return false
	}
}

func Sign(appkey string, appSecret string, random int, timestamp int64, hash crypto.Hash) (string, error) {
	a := appkey + appSecret + strconv.Itoa(random) + strconv.Itoa(int(timestamp))
	h := hash.New()
	io.WriteString(h, a)
	test := h.Sum(nil)
	sign := fmt.Sprintf("%x", test)
	return sign, nil
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
