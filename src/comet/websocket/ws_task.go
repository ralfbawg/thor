package websocket

import (
	"common/logging"
	ws "github.com/gorilla/websocket"
	"github.com/panjf2000/ants"
	"sync/atomic"
	"time"
	"util"
	"context"
	"fmt"
)

const (
	broadcastBatchNum = 2
)

var (
	wrPoolExtendFactor = 0.8
	wrPoolDefaultSize  = 10000
	wrPool, _          = ants.NewPool(wrPoolDefaultSize)
)

type WsTask struct {
	app *WsApp
	// App id
	appId string
	//clients *util.ConcurrentMap

	clients util.ConcMap

	index int

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *WsTaskClient

	// Unregister requests from clients.
	unregister chan *WsTaskClient

	clientCount int64

	incr chan int64
}
type Ws_context struct {
	//context内容
	context.Context
}

func (task *WsTask) AddClient(uid string, conn *ws.Conn) *WsTaskClient {
	if task == nil || task.clients == nil {
		logging.Info(util.AOrB(func() bool {
			return task == nil
		}, "task is empty", "task clients is empty").(string))
		return nil
	}
	if client, ok := task.clients.Get(uid); ok && client != nil { //TODO 如果存在，如何处理,暂时先断开，删除
		task.unregister <- client.(*WsTaskClient)
		//task.clientsIndex[id].conn.Close()
		//delete(task.clients, task.clientsIndex[id])
		//delete(task.clientsIndex, id)
	}

	client := &WsTaskClient{
		task: task,
		conn: conn,
		uid:  uid,
		send: make(chan []byte, 10),
	}
	client.task.register <- client

	util.SubmitTaskAndResize(wrPool, wrPoolDefaultSize, wrPoolExtendFactor, client.readGoroutine, client.writeGoroutine) //fixme funcs有可能有同步问题
	//client.Send([]byte(hiMesaage + "," + client.uid))                                                                                       //fixme 第一次连接发送，方便测试
	client.Send([]byte(fmt.Sprintf(hiMesaageJson, client.uid)))
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
func (task *WsTask) GetClients() map[string]interface{} {
	return task.clients.Items()
}
func (task *WsTask) GetAppId() string {
	return task.appId
}
func (task *WsTask) Init() string {
	return task.appId
}

func NewWsTask(app *WsApp, index int) *WsTask {
	task := &WsTask{
		appId: app.appId,
		app:   app,
		index: index,
		//clients:      make(map[*WsTaskClient]bool),
		clients: util.NewConcMap(),
		//clientsIndex: util.NewConcurrentMap(),
		broadcast:  make(chan []byte, 10),
		register:   make(chan *WsTaskClient, 2000),
		unregister: make(chan *WsTaskClient, 2000),
		incr:       make(chan int64, 2000),
	}
	manager.register <- app
	ants.Submit(task.Run)
	ants.Submit(task.statistic)
	//go task.Run()
	//go task.statistic()
	return task
}

func (task *WsTask) Run() {
	defer func() {
		recover()
	}()
	for {
		select {
		case client := <-task.register:
			task.clients.Set(client.uid, client)
			task.incr <- 1
		case client := <-task.unregister:
			if tClient, ok := task.clients.Get(client.uid); ok && tClient != nil {
				task.clients.Remove(client.uid)
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
			atomic.AddInt64(&task.clientCount, in)
			if in < 0 && count <= 0 {
				atomic.StoreInt64(&task.app.clientCount, 0)
				//if count < 0 {
				//	atomic.StoreInt64(&task.wsManager.TaskCount,0)
				//	atomic.AddInt64(&task.wsManager.TaskCount, int64(1-atomic.LoadInt64(&task.wsManager.TaskCount)))
				//} else {
				//	atomic.AddInt64(&task.wsManager.TaskCount, int64(1-atomic.LoadInt64(&task.wsManager.TaskCount)))
				//	manager.tasks.Remove(task.appId)
				//}
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
