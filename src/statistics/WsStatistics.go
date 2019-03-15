package statistics

import "websocket"


func GetTaskCount() int {
	return websocket.GetWsManager().TaskCount
}
func GetClientCount(appId string) int {
	return websocket.GetWsManager().GetOrCreateTask(appId).GetClientCount()
}