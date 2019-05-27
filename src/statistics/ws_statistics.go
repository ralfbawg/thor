package statistics

import "comet/websocket"

func GetTaskCount() int64 {
	return websocket.GetWsManager().GetTaskCount()
}
//func GetClientCount(appId string) int64 {
//	return websocket.GetWsManager().GetOrCreateTask(appId).GetClientCount()
//}
func GetAllClientCount() int64 {
	return websocket.GetWsManager().GetAllClientCount()
}
