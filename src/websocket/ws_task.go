package websocket

import (
	ws "github.com/gorilla/websocket"
	"util"
)

type WsTask struct {
	wsManager *WsManager
	// App id
	appId string
	// Registered clients.
	//clients *util.ConcurrentMap

	clients *util.ConcurrentMap

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *WsTaskClient

	// Unregister requests from clients.
	unregister chan *WsTaskClient

	clientCount int
}

func (task *WsTask) AddClient(id string, conn *ws.Conn) *WsTaskClient {
	if client := task.clients.Get(id); client != nil { //TODO 如果存在，如何处理,暂时先断开，删除
		task.unregister <- client.(*WsTaskClient)
		//task.clientsIndex[id].conn.Close()
		//delete(task.clients, task.clientsIndex[id])
		//delete(task.clientsIndex, id)
	}

	client := &WsTaskClient{
		task: task,
		conn: conn,
		id:   id,
		send: make(chan []byte),
	}
	client.task.register <- client
	go client.readGoroutine()
	go client.writeGoroutine()
	return client
}

func (task *WsTask) Broadcast(msg []byte) {
	for _, v := range task.clients.Map {
		v.(*WsTaskClient).send <- msg
	}
}
func (task *WsTask) GetClientCount() int {
	return task.clientCount
}

func NewWsTask(appId string, manager *WsManager) *WsTask {
	task := &WsTask{
		appId:     appId,
		wsManager: manager,
		//clients:      make(map[*WsTaskClient]bool),
		clients: util.NewConcurrentMap(),
		//clientsIndex: util.NewConcurrentMap(),
		broadcast:  make(chan []byte),
		register:   make(chan *WsTaskClient, 1000),
		unregister: make(chan *WsTaskClient, 1000),
	}
	go task.Run()
	return task
}

func (task *WsTask) Run() {
	for {
		select {
		case client := <-task.register:
			task.clients.Put(client.id, client)
			task.clientCount++
		case client := <-task.unregister:
			if tClient := task.clients.Get(client.id); tClient != nil {
				task.clients.Del(client.id)
				//task.clientsIndex.Del(client.id)
				//delete(task.clientsIndex, client.id)
				close(client.send)
				task.clientCount--
			}
		case message := <-task.broadcast:
			task.clients.Foreach(func(k string, v interface{}) {
				v.(*WsTaskClient).send <- message
			})
			//for client := range task.clients {
			//	client.send <- message
			//}
		}
	}
}
