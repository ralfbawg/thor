package task

import (
	"net/url"
	"util/uuid"
	"strconv"
)

const (
	appIdParam  = "appId"
	taskIdParam = "taskId"
	appKeyParam = "appKey"
	uidParam    = "uid"
)

/*
 验证app信息
*/
func VerifyAppInfo(param url.Values) (string, int, string, bool) {
	appId := param.Get(appIdParam)
	appKey := param.Get(appKeyParam)
	uid := param.Get(uidParam)
	taskId, err := strconv.Atoi(param.Get(taskIdParam))
	if err != nil {
		taskId = 0
	}
	if uid == "" {
		uid = uuid.Generate().String()
	}
	//logging.Debug("app id is %s,app key is %s,uid is %s", appId, appKey, id)
	return VerifyAppInfo2(appId, taskId, uid, appKey)
}

/*
 验证app信息
*/
func VerifyAppInfo2(appId string, taskId int, uid string, appkey string) (string, int, string, bool) {
	if uid == "" {
		uid = uuid.Generate().String()
	}
	//logging.Debug("app id is %s,app key is %s,uid is %s", appId, appKey, id)
	//TODO 通过db查询确认
	//return id, appKey != "fffasdfasdf" && id != "asdfasdfasd"
	return appId, taskId, uid, true
}
