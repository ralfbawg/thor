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

type DiagnoseStat struct {
	Alloc    float64 `json:"alloc"`
	Inuse    float64 `json:"inuse"`
	Idle     float64 `json:"idle"`
	Sys      float64 `json:"sys"`
	Released float64 `json:"released"`
}

type CpuStat struct {
	Usage float64
	Busy  float64
	Total float64
}
