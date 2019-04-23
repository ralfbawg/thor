package task

import (
	"db"
	"net/url"
	"task/provider"
	"util/uuid"
	"strconv"
)

const (
	appIdParam  = "appId"
	taskIdParam = "taskId"
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
func VerifyAppInfo(param url.Values) (string, int, string, bool) {
	appId := param.Get(appIdParam)
	//appKey := param.Get(appKeyParam)
	id := param.Get(uidParam)
	taskId, err := strconv.Atoi(param.Get(taskIdParam))
	if err != nil {
		taskId = 0
	}
	if id == "" {
		id = uuid.Generate().String()
	}
	//logging.Debug("app id is %s,app key is %s,uid is %s", appId, appKey, id)
	//TODO 通过db查询确认
	//return id, appKey != "fffasdfasdf" && id != "asdfasdfasd"
	return appId, taskId, id, true
}
