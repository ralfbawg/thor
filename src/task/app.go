package task

import (
	"db"
	"net/url"
	"task/provider"
	"util/uuid"
)

const (
	appIdParam  = "appId"
	appKeyParam = "appKey"
	uidParam    = "uid"
)

type AppTask struct {
	Provider provider.BaseProvider
}

func NewAppTask() *AppTask {
	return nil
}

func InitAppTask() {
	db.GetAppDb()
}

/*
 验证app信息
*/
func VerifyAppInfo(param url.Values) (string, string, bool) {
	appId := param.Get(appIdParam)
	//appKey := param.Get(appKeyParam)
	id := param.Get(uidParam)
	if id == "" {
		id = uuid.Generate().String()
	}
	//logging.Debug("app id is %s,app key is %s,uid is %s", appId, appKey, id)
	//TODO 通过db查询确认
	//return id, appKey != "fffasdfasdf" && id != "asdfasdfasd"
	return appId, id, true
}
