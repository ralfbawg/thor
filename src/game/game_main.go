package game

import (
	"common/logging"
	"errors"
	ws "github.com/gorilla/websocket"
	"net/http"
	"sync/atomic"
	"time"
)

var (
	gameRoomsArr = make([]uint32, 100)
	posArr       = [...]uint32{
		1 << 31,
		1 << 30,
		1 << 29,
		1 << 28,
		1 << 27,
		1 << 26,
		1 << 25,
		1 << 24,
		1 << 23,
		1 << 22,
		1 << 21,
		1 << 20,
		1 << 19,
		1 << 18,
		1 << 17,
		1 << 16,
		1 << 15,
		1 << 14,
		1 << 13,
		1 << 12,
		1 << 11,
		1 << 10,
		1 << 9,
		1 << 8,
		1 << 7,
		1 << 6,
		1 << 5,
		1 << 4,
		1 << 3,
		1 << 2,
		1 << 1,
		1}
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
)

type GameRoom struct {
	clientA     *GameClient
	clientB     *GameClient
	broadcast   chan []byte
	npc         func(clientA *GameClient, clientB *GameClient)
	gameRunning bool
	gameStart   time.Duration
	unregister  chan *GameClient
}

func (gr *GameRoom) AddClient(client *GameClient) bool {
	if gr.clientA == nil {
		gr.clientA = client
		return true
	} else if gr.clientB != nil {
		gr.clientB = client
		return true
	}
	return false
}

type GameRooms []*GameRoom

func NewGameRooms() GameRooms {
	tmp := make(GameRooms, 1000)
	for i := 0; i < 1000; i++ {
		tmp[i] = &GameRoom{
			broadcast: make(chan []byte),
			//npc: func(clientA *GameClient, clientB *GameClient) {
			//	if clientA != nil && clientB! = nil && gameRunning{
			//
			//	}
			//},
		}
	}
	GameRoomsInst = tmp
	return tmp
}

func GameDispatch(w http.ResponseWriter, r *http.Request) {
	//param := r.URL.Query()
	//logging.Debug(param.Get("appId"))

	if conn, err := upgrader.Upgrade(w, r, nil); err != nil {
		logging.Error("哦活,error:%s", err)
	} else {
		//task := manager.GetOrCreateTask(appId)
		//task.AddClient(id, conn)
	}

}

func CreateOrGetGameRoom() (*GameRoom, error) {
	if id, err := CreateOrGetGameRoomId(); err != nil {
		return nil, err
	} else {
		return GameRoomsInst[id], nil
	}
}

func CreateOrGetGameRoomId() (int, error) {
	for i := 0; i < len(gameRoomsArr); i++ {
		j := gameRoomsArr[i]
		for m := 0; m < len(posArr); m = m + 2 {
			n := posArr[m]
			o := posArr[m+1]
			if j&n != n {
				tmp := j ^ n
				if atomic.CompareAndSwapUint32(&j, j, tmp) {
					return (i*32 + m + 1) / 2, nil
				}
				break
			} else if j&o != o {
				tmp := j ^ o
				if atomic.CompareAndSwapUint32(&j, j, tmp) {
					return (i*32 + m + 1 + 1) / 2, nil
				}
				break
			}
		}
	}
	return 0, errors.New("满了")
}
