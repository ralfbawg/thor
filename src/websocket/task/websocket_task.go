package task

type WsTask struct {
	// App id
	AppId string
	// Registered clients.
	clients map[*WsTaskClient]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *WsTaskClient

	// Unregister requests from clients.
	unregister chan *WsTaskClient
}

func (task *WsTask)AddClient(string)  {

}