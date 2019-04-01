package game

import (
	"net/http"
	"task"
	"common/logging"
)

type game struct {

}

func (g *game)CreateOrFindEmptyRoom()  {
	
}


type gameRoom struct {
	
}

func GameDispatch(w http.ResponseWriter, r *http.Request) {
	param := r.URL.Query()
	//logging.Debug(param.Get("appId"))
	if appId, id, exist := VerifyAppInfo(param); exist == true {
		if conn, err := upgrade.Upgrade(w, r, nil); err != nil {
			logging.Error("哦活,error:%s", err)
		} else {
			//task := manager.GetOrCreateTask(appId)
			//task.AddClient(id, conn)
		}
	} else {
		w.Write([]byte("appId 错误或者不存在"))
	}
}
