package api

type ListOnlineUsersResp struct {
	Code int
	List []*ClientInfo
}
type ClientInfo struct {
	AppId    string //业务id
	ClientId string //客户端id
	ClientIp string //客户端ip
}

type ConnectingUsersResp struct {
	Code  int
	Users int64
}
