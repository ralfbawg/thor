package websocket

import (
	ws "github.com/gorilla/websocket"
	"util"
	"sync/atomic"
	"common/logging"
	"time"
	"github.com/panjf2000/ants"
)

const (
	broadcastBatchNum = 2
)

var (
	wrPoolExtendFactor = 0.8
	wrPoolDefaultSize  = 10000
	wrPool, _          = ants.NewPool(wrPoolDefaultSize)
	funcs              = make([]func(), 0)
)

type WsTask struct {
	wsManager *WsManager
	// App id
	appId string
	// Registered clients.
	//clients *util.ConcurrentMap

	clients util.ConcMap

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
	if task == nil {
		logging.Info("task is empty")
	}
	if task.clients == nil {
		logging.Info("task clients is empty")
	}
	if client, ok := task.clients.Get(id); ok && client != nil { //TODO 如果存在，如何处理,暂时先断开，删除
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

	submitTaskAndResize(wrPool, append(funcs[:0], client.readGoroutine, client.writeGoroutine))
	//wrPool.Submit(client.readGoroutine)
	//go client.readGoroutine()
	//go client.writeGoroutine()
	client.Send([]byte(hiMesaage + "," + client.id)) //fixme 第一次连接发送，方便测试
	return client
}

func (task *WsTask) Broadcast(msg []byte) {
	start := time.Now()
	//modInt := 0
	//if task.clientCount%broadcastBatchNum > 0 {
	//	modInt = 1
	//}
	//batchCount := task.clientCount/broadcastBatchNum + int64(modInt)
	//timeTask(1*time.Second, int(batchCount), func(param ...interface{}) { //FIXME 需要解决map的乱序问题
	//	n := param[0].(int64)
	//	tmpCount := 0
	//	task.clients.Keys()
	//	task.clients.IterCb(func(key string, v interface{}) {
	//		if tmpCount < (n-1*(broadcastBatchNum)) || tmpCount >= int(n*broadcastBatchNum) {
	//			return
	//		}
	//		if v != nil {
	//			v.(*WsTaskClient).send <- msg
	//		}
	//		tmpCount++
	//	})
	//})
	task.clients.IterCb(func(key string, v interface{}) {
		if v != nil {
			v.(*WsTaskClient).send <- msg
		}
	})
	end := time.Now()
	logging.Debug("broadcast time cost %f second", end.Sub(start).Seconds())
}
func (task *WsTask) GetClientCount() int64 {
	return atomic.LoadInt64(&task.clientCount)
}
func (task *WsTask) GetAppId() string {
	return task.appId
}
func (task *WsTask) Init() string {
	return task.appId
}


func NewWsTask(appId string, manager *WsManager) *WsTask {
	task := &WsTask{
		appId:     appId,
		wsManager: manager,
		//clients:      make(map[*WsTaskClient]bool),
		clients: util.NewConcMap(),
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
			task.clients.Set(client.id, client)
			task.incr <- 1
		case client := <-task.unregister:
			if tClient, ok := task.clients.Get(client.id); ok && tClient != nil {
				task.clients.Remove(client.id)
				close(client.send)
				task.incr <- -1
			}
		case message := <-task.broadcast:
			task.clients.IterCb(func(k string, v interface{}) {
				v.(*WsTaskClient).send <- message
			})
		}
	}
}
func (task *WsTask) GetClient(uid string) *WsTaskClient {
	if client, ok := task.clients.Get(uid); ok && client != nil {
		return client.(*WsTaskClient)
	} else {
		return nil
	}
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
			if in < 0 && count <= 0 {
				if count < 0 {
					atomic.AddInt64(&task.wsManager.TaskCount, int64(1-task.wsManager.TaskCount))
				} else {
					atomic.AddInt64(&task.wsManager.TaskCount, int64(1-task.wsManager.TaskCount))
					manager.tasks.Remove(task.appId)
				}

			}
		}
	}
}

//分批任务,
// t 任务间隔
//count 任务执行次数
//f 任务本身
func timeTask(t time.Duration, count int, f func(param ...interface{})) {
	go func() {
		tick := time.NewTicker(t)
		initCount := 0
		for range tick.C {
			if initCount < count {
				f([...]interface{}{initCount})
				initCount++
			} else {
				tick.Stop()
			}
		}
	}()
}

func submitTaskAndResize(pool *ants.Pool, f []func()) {
	if float64(pool.Running())/float64(wrPoolDefaultSize) > wrPoolExtendFactor {
		wrPoolDefaultSize *= 2
		pool.Tune(wrPoolDefaultSize)
	}
	for _, v := range f {
		pool.Submit(v)
	}
}
