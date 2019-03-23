package websocket

import (
	ws "github.com/gorilla/websocket"
	"util"
	"sync/atomic"
	"time"
	"common/logging"
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

	incr chan int64
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
		send: make(chan []byte, 10),
	}
	client.task.register <- client
	go client.readGoroutine()
	go client.writeGoroutine()
	client.Send([]byte(hiMesaage + "," + client.id)) //fixme 第一次连接发送，方便测试
	return client
}

func (task *WsTask) Broadcast(msg []byte) {
	start := time.Now()
	task.clients.Foreach(func(s string, i interface{}) {
		i.(*WsTaskClient).send <- msg
	})
	end := time.Now()
	logging.Debug("broadcast time cost %d second", end.Sub(start).Seconds())
}
func (task *WsTask) GetClientCount() int64 {
	return atomic.LoadInt64(&task.clientCount)
}
func (task *WsTask) GetAppId() string {
	return task.appId
}

func NewWsTask(appId string, manager *WsManager) *WsTask {
	task := &WsTask{
		appId:     appId,
		wsManager: manager,
		//clients:      make(map[*WsTaskClient]bool),
		clients: util.NewConcurrentMap(),
		//clientsIndex: util.NewConcurrentMap(),
		broadcast:  make(chan []byte, 10),
		register:   make(chan *WsTaskClient, 2000),
		unregister: make(chan *WsTaskClient, 2000),
		incr:       make(chan int64, 2000),
	}
	manager.register <- task
	go task.Run()
	go task.statistic()
	return task
}

func (task *WsTask) Run() {
	defer func() {
		recover()
	}()
	for {
		select {
		case client := <-task.register:
			task.clients.Put(client.id, client)
			task.incr <- 1
		case client := <-task.unregister:
			if tClient := task.clients.Get(client.id); tClient != nil {
				task.clients.Del(client.id)
				close(client.send)
				task.incr <- -1
			}
		case message := <-task.broadcast:
			task.clients.Foreach(func(k string, v interface{}) {
				v.(*WsTaskClient).send <- message
			})
		}
	}
}
func (task *WsTask) GetClient(uid string) *WsTaskClient {
	return task.clients.Get(uid).(*WsTaskClient)
}
func (task *WsTask) statistic() {
	for {
		select {
		case in := <-task.incr:
			intLen := len(task.incr)
			for i := 0; i < intLen; i++ {
				t := <-task.incr
				in = in + t
			}
			count := atomic.AddInt64(&task.clientCount, in)
			atomic.AddInt64(&task.wsManager.ClientCount, in)
			if in < 0 && count == 0 {
				atomic.AddInt64(&task.wsManager.TaskCount, in)
			}
		}
	}
}
