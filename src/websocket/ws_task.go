package websocket

import (
	ws "github.com/gorilla/websocket"
	"util"
	"sync/atomic"
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

	clientCount int64
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
func (task *WsTask) GetClientCount() int64 {
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
	manager.register <- task
	go task.Run()
	return task
}

func (task *WsTask) Run() {
	for {
		select {
		case client := <-task.register:
			task.clients.Put(client.id, client)
			atomic.AddInt64(&task.clientCount, 1)
			atomic.AddInt64(&task.wsManager.ClientCount, 1)
		case client := <-task.unregister:
			if tClient := task.clients.Get(client.id); tClient != nil {
				task.clients.Del(client.id)
				close(client.send)
				atomic.AddInt64(&task.clientCount, -1)
				atomic.AddInt64(&task.wsManager.ClientCount, -1)
			}
		case message := <-task.broadcast:
			task.clients.Foreach(func(k string, v interface{}) {
				v.(*WsTaskClient).send <- message
			})
		}
	}
}
