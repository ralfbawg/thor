package rpc

import (
	"comet/tcp"
)

var RpcClient = Init()

type RouterServiceClient struct {
}

func Init() *RouterServiceClient {
	return &RouterServiceClient{

	}
}
func (s *RouterServiceClient) SendMsg(appId string, taskId int, uid string, msg []byte) error {
	tcp.SendMsg(appId, taskId, uid, msg) //TODO 切换成rpc
	//d := client.NewPeer2PeerDiscovery("tcp@"+addr, "")
	//xclient := client.NewXClient("Arith", client.Failtry, client.RandomSelect, d, client.DefaultOption)
	//defer xclient.Close()
	//err := xclient.Call(context.Background(), "Mul", args, reply, nil)
	return nil
}
func (client *RouterServiceClient) Register(s string, test func()) {

}
