package websocket

import (
	ws "github.com/gorilla/websocket"
)

type WsTask struct {
	wsManager *WsManager
	// App id
	appId string
	// Registered clients.
	clients map[*WsTaskClient]bool

	clientsIndex map[string]*WsTaskClient

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *WsTaskClient

	// Unregister requests from clients.
	unregister chan *WsTaskClient
}

func (task *WsTask) AddClient(id string, conn *ws.Conn) *WsTaskClient {
	if client := task.clientsIndex[id]; client != nil { //TODO 如果存在，如何处理,暂时先断开，删除
		task.unregister <- client
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
	for k, _ := range task.clients {
		k.send <- msg
	}
}

func NewWsTask(appId string, manager *WsManager) *WsTask {
	task := &WsTask{
		appId:        appId,
		wsManager:    manager,
		clients:      make(map[*WsTaskClient]bool),
		clientsIndex: make(map[string]*WsTaskClient),
		broadcast:    make(chan []byte),
		register:     make(chan *WsTaskClient),
		unregister:   make(chan *WsTaskClient),
	}
	go task.Run()
	return task
}

func (task *WsTask) Run() {
	for {
		select {
		case client := <-task.register:
			task.clients[client] = true
			task.clientsIndex[client.id] = client
		case client := <-task.unregister:
			if _, ok := task.clients[client]; ok {
				delete(task.clients, client)
				delete(task.clientsIndex, client.id)
				close(client.send)
			}
		case message := <-task.broadcast:
			for client := range task.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(task.clients, client)
				}
			}
		}
	}
}
