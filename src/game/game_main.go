package game

import (
	ws "github.com/gorilla/websocket"
	"net/http"
	"sync/atomic"
	"time"
	"common/logging"
	"github.com/panjf2000/ants"
	"strconv"
)

const (
	GAME_STATUS_PREPARE = iota
	GAME_STATUS_READY
	GAME_STATUS_RUNNING
	GAME_STATUS_FINISH
	GAME_STATUS_EMPTY
	GAME_STATUS         = iota
)

var (
	GameRoomsArr = make([]uint32, 100000)
	//32的异或数组，从高位开始
	posArr   = [...]uint32{1 << 31, 1 << 30, 1 << 29, 1 << 28, 1 << 27, 1 << 26, 1 << 25, 1 << 24, 1 << 23, 1 << 22, 1 << 21, 1 << 20, 1 << 19, 1 << 18, 1 << 17, 1 << 16, 1 << 15, 1 << 14, 1 << 13, 1 << 12, 1 << 11, 1 << 10, 1 << 9, 1 << 8, 1 << 7, 1 << 6, 1 << 5, 1 << 4, 1 << 3, 1 << 2, 1 << 1, 1}
	upgrader = ws.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		// 允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	GameRoomsInst = func() GameRooms {
		return NewGameRooms()
	}()
	uintfullArr = [...]uint32{uint32(65535), uint32(255), uint32(15), uint32(3), uint32(1)}
)

type GameRoom struct {
	clientA    *GameClient
	clientB    *GameClient
	obs        []*GameClient
	broadcast  chan []byte
	npc        func(clientA *GameClient, clientB *GameClient, running bool)
	unregister chan *GameClient
	status     int8
	gameStatus int8
	statusC    chan int8
}

func (gr *GameRoom) AddClient(client *GameClient) bool {
	if gr.status != GAME_STATUS_PREPARE {
		return false
	}
	if gr.clientA == nil {
		gr.clientA = client
		return true
	} else if gr.clientB != nil {
		gr.clientB = client
		return true
	}
	return false
}
func (gr *GameRoom) Run() {
	for {
		select {
		case s := <-gr.statusC:
			switch s {
			case GAME_STATUS_READY:
				gr.start = (time.Now().Add(5 * time.Second))
			}
			ants.Submit(gr.RunGame)
			break
		}
	}

}
func (gr *GameRoom) RunGame() {

	for {
		select {
		case <-gr.ready:

		}
	}
	for {
		select {
		case msg := <-gr.broadcast:
			gr.clientA.Send(msg)
			gr.clientB.Send(msg)
		}
	}
}

type GameRooms []*GameRoom

func NewGameRooms() GameRooms {
	tmp := make(GameRooms, 1000)
	for i := 0; i < 1000; i++ {
		tmp[i] = &GameRoom{
			broadcast: make(chan []byte, 10),
			npc: func(clientA *GameClient, clientB *GameClient, running bool) {
			},
		}
	}
	return tmp
}

func GameDispatch(w http.ResponseWriter, r *http.Request) {
	//param := r.URL.Query()
	//logging.Debug(param.Get("appId"))

	//if conn, err := upgrader.Upgrade(w, r, nil); err != nil {
	//	logging.Error("哦活,error:%s", err)
	//} else {
	//	//task := manager.GetOrCreateTask(appId)
	//	//task.AddClient(id, conn)
	//}

}

func
CreateOrGetGameRoom() (*GameRoom, error) {
	if id, err := CreateOrGetGameRoomId(); err != nil {
		return nil, err
	} else {
		return GameRoomsInst[id], nil
	}
}

func
CreateOrGetGameRoomId() (int, error) {
	for i := 0; i < len(GameRoomsArr); i++ {
		for j := atomic.LoadUint32(&GameRoomsArr[i]); j < ^uint32(0); {
			if n, b := findRoomByBinarySearch(&GameRoomsArr[i], j, 0, 32, 0); b {
				return n + i*16, nil
			}
		}
		//gameRoomNo := int(findRoomByBinarySearch(j, 0, 32, 0))

		//for m := 0; m < len(posArr); m = m + 2 { //todo 顺序查找,未来估计要改成二分查找
		//	n := posArr[m]
		//	o := posArr[m+1]
		//	if j&n != n {
		//		tmp := j ^ n
		//		if atomic.CompareAndSwapUint32(&j, j, tmp) {
		//			return (i*32 + m + 1) / 2, nil
		//		}
		//		break
		//	} else if j&o != o {
		//		tmp := j ^ o
		//		if atomic.CompareAndSwapUint32(&j, j, tmp) {
		//			return (i*32 + m + 1 + 1) / 2, nil
		//		}
		//		break
		//	}
		//}
	}
	logging.Debug("我没有找到房间")
	return CreateOrGetGameRoomId()
	//return 0, errors.New("满了")
}

func findRoomByBinarySearch(ptr *uint32, p uint32, start uint32, end uint32, uintArrIndex int) (int, bool) {
	tValue := atomic.LoadUint32(ptr)
	if p == ^uint32(0) { //满了
		return -1, false
	}
	if end-start <= 1 {
		gameRoomNo := (start+1)/2 + (start+1)%2
		n := posArr[start]
		tmp := tValue ^ n
		return int(gameRoomNo), atomic.CompareAndSwapUint32(ptr, tValue, tmp)
	}
	offset := (end - start) / 2
	tmp := uintfullArr[uintArrIndex]
	if tmp&(p>>(offset+(32-end))) < tmp {
		return findRoomByBinarySearch(ptr, p, start, end-offset, uintArrIndex+1) //左,会右移，低位不需要置零
	} else {
		return findRoomByBinarySearch(ptr, p&(^(tmp << (offset + 32 - end))), start+offset, end, uintArrIndex+1) //右，高位需要置零
	}

}
func
FindNotEmptyRoom() {
	for k, v := range GameRoomsArr {
		//if (v != ^uint32(0) && v != 0) || (k%8 == 0 && v != 0) {
		//	logging.Info("the room(%d) is not full,value=%d", k+1, v)
		//}
		if (v != ^uint32(0) && v != 0) {
			logging.Info("the room(%d) is not full,value=%d", k+1, v)
		}
	}
}

type Game struct {
	gr      *GameRoom
	statusC chan int8
	start   time.Duration
	bloodA  int8
	bloodB  int8
	clientA *GameClient
	clientB *GameClient
}

func (g *Game) Run() {
	ALastHit := time.Now()
	npcLast := time.Now()
	BLastHit := time.Now()
	for {
		select {
		case status := <-g.statusC:

		case rmsg := <-g.clientA.read:
			switch string(rmsg) {
			case EA:
				if time.Now().Sub(npcLast) > 5*time.Second {

				}
			}


		case rmsg := <-g.clientB.read:
			switch string(rmsg) {
			case EA:
				if time.Now().Sub(npcLast) > 5*time.Second {
					g.gr.broadcast <- []byte(OB2 + "A")
				} else {
					g.gr.broadcast <- []byte(OB2 + "B")
				}
			}
		}
	}

}
