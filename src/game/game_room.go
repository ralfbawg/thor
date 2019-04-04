package game

import (
	"github.com/panjf2000/ants"
	"time"
	"sync/atomic"
	"common/logging"
	"errors"
)

const (
	GAME_STATUS_PREPARE = iota
	GAME_STATUS_READY
	GAME_STATUS_RUNNING
	GAME_STATUS_FINISH
	GAME_STATUS_EMPTY
	GAME_ERROR_FIND     = "EF"
	GAME_EVENT_START3   = "GS3"
	GAME_EVENT_START5   = "GS5"
)

var (
	GameRoomsArr = make([]uint32, 1000)
	//32的异或数组，从高位开始
	posArr      = [...]uint32{1 << 31, 1 << 30, 1 << 29, 1 << 28, 1 << 27, 1 << 26, 1 << 25, 1 << 24, 1 << 23, 1 << 22, 1 << 21, 1 << 20, 1 << 19, 1 << 18, 1 << 17, 1 << 16, 1 << 15, 1 << 14, 1 << 13, 1 << 12, 1 << 11, 1 << 10, 1 << 9, 1 << 8, 1 << 7, 1 << 6, 1 << 5, 1 << 4, 1 << 3, 1 << 2, 1 << 1, 1}
	uintfullArr = [...]uint32{uint32(65535), uint32(255), uint32(15), uint32(3), uint32(1)}
)

type GameRoom struct {
	index      int
	clientA    *GameClient
	clientB    *GameClient
	gm         *GameMall
	obs        []*GameClient
	broadcast  chan []byte
	npc        func(clientA *GameClient, clientB *GameClient, running bool)
	status     int8
	gameStatus int8
	statusC    chan int8
}

type GameRooms []*GameRoom

func (gr *GameRoom) ExitClient(client *GameClient) bool {
	if gr.clientA.id == client.id {
		gr.gm.waitingClients.Set(client.id, client)
		gr.clientA = nil
	}
	if gr.clientB.id == client.id {
		gr.gm.waitingClients.Set(client.id, client)
		gr.clientB = nil
	}
	if gr.clientA == nil && gr.clientB == nil {
		gr.statusC <- GAME_STATUS_EMPTY
	}
	return true
}

func (gr *GameRoom) AddClient(client *GameClient) bool {
	if gr.status != GAME_STATUS_PREPARE {
		return false
	}
	if gr.clientA == nil {
		gr.clientA = client
		client.pos = ROOM_POS_A
		if gr.clientA != nil && gr.clientB != nil {
			gr.statusC <- GAME_STATUS_READY

		}
		return true
	} else if gr.clientB != nil {
		gr.clientB = client
		client.pos = ROOM_POS_B
		if gr.clientA != nil && gr.clientB != nil {
			gr.statusC <- GAME_STATUS_READY

		}
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
				gr.status = GAME_STATUS_READY
				//gr.start = (time.Now().Add(5 * time.Second))
				gr.BroadCast([]byte(GAME_EVENT_START3))
				ants.Submit(func() { //开始
					time.AfterFunc(3*time.Second, func() {
						gr.statusC <- GAME_STATUS_RUNNING
					})
				})
				ants.Submit(func() { //结束
					time.AfterFunc(33*time.Second, func() {
						gr.statusC <- GAME_STATUS_FINISH
					})
				})
			case GAME_STATUS_EMPTY:
				if ResetRoomStatus(gr.index) {
					gr.status = GAME_STATUS_PREPARE
				}
			}
			ants.Submit(gr.RunGame)
			break
		}
	}

}

func (gr *GameRoom) RunGame() {

}
func (room *GameRoom) BroadCast(msg []byte) {
	room.clientB.Send(msg)
	room.clientA.Send(msg)
}

//func (room *GameRoom) Remove(id string) {
//	if room.clientA.id == id {
//
//	}
//}

func NewGameRooms() GameRooms {
	tmp := make(GameRooms, 1000)
	for i := 0; i < 1000; i++ {
		tmp[i] = &GameRoom{
			index:     i,
			broadcast: make(chan []byte, 10),
			status:    GAME_STATUS_PREPARE,
			statusC:   make(chan int8, 2),
			gm:        GameMallInst,
			npc: func(clientA *GameClient, clientB *GameClient, running bool) {
			},
		}
	}
	return tmp
}
func (gr GameRooms) CreateOrGetGameRoom() (*GameRoom, error) {
	if id, err := CreateOrGetGameRoomId(); err != nil {
		return nil, err
	} else {
		return gr[id], nil
	}
}

func CreateOrGetGameRoomId() (int, error) {
	for i := 0; i < len(GameRoomsArr); i++ {
		for j := atomic.LoadUint32(&GameRoomsArr[i]); j < ^uint32(0); {
			if n, b := findRoomByBinarySearch(&GameRoomsArr[i], j, 0, 32, 0); b {
				return n + i*16, nil
			}
		}
	}
	return 0, errors.New("满了")
}

//重置房间状态
func ResetRoomStatus(index int) bool {
	i := index / 16
	j := 0
	if index%16 != 0 {
		j = 1
	}
	old := GameRoomsArr[i-1+j]
	new := uint32(3)
	if index%16 != 0 {
		new = uint32(3) << uint32((16-index%16)*2)
	}
	new = old & new
	return atomic.CompareAndSwapUint32(&GameRoomsArr[i-1+j], old, new)
}

// fixme 二分查找，有一点点的同步的问题
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

// 顺序查找
func findRoomByNormal() (int, bool) {
	for i := 0; i < len(GameRoomsArr); i++ {
		j := GameRoomsArr[i]
		for m := 0; m < len(posArr); m = m + 2 {
			n := posArr[m]
			o := posArr[m+1]
			if j&n != n {
				tmp := j ^ n
				if atomic.CompareAndSwapUint32(&j, j, tmp) {
					return (i*32 + m + 1) / 2, true
				}
				break
			} else if j&o != o {
				tmp := j ^ o
				if atomic.CompareAndSwapUint32(&j, j, tmp) {
					return (i*32 + m + 1 + 1) / 2, true
				}
				break
			}
		}
	}
	return -1, false
}

func FindNotEmptyRoom() {
	for k, v := range GameRoomsArr {
		if (v != ^uint32(0) && v != 0) {
			logging.Info("the room(%d) is not full,value=%d", k+1, v)
		}
	}
}
