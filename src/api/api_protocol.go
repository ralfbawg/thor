package api

type ListOnlineUsersResp struct {
	Code int
	List []*ClientInfo
}
type ClientInfo struct {
	ClientId string //客户端id
	ClientIp string //客户端ip
	Status   int    //客户端状态, 0:链接中 1:匹配中 2:游戏中
}

type ConnectingUsersResp struct {
	Code  int
	Users int
}
