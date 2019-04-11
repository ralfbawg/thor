package game

import (
	"common/logging"
	"encoding/json"
	"github.com/panjf2000/ants"
	"math/rand"
	"sync/atomic"
	"time"
)

const (
	EA  = "EA"  //攻击事件
	EAA = "EAA" //A攻击事件
	EAB = "EAB" //B攻击事件
	EN1 = "EN1" //npc事件1
	EN2 = "EN2" //npc事件2
	OB1 = "OB1" //掉血1
	OB2 = "OB2" //掉血2
)

type GameI interface {
	RunGame(gr *GameRoom)
	Init(gr *GameRoom)
	OnEvent(gameRoomStatus int)
}

type BaseGame struct {
	name       string
	code       string
	prefixTime time.Duration
	duringTime time.Duration
	npc        func(pos int32) interface{}
	startTime  time.Time
	gameRoom   *GameRoom
	GameI
}

func (bg *BaseGame) RunGame(gr *GameRoom) {
	logging.Info("default RuaGame")
}
func (bg *BaseGame) Init(gr *GameRoom) {
	logging.Info("default Init")
}
func (bg *BaseGame) timeCounter(gr *GameRoom) {
	ants.Submit(func() { //开始
		time.AfterFunc(bg.prefixTime, func() {
			gr.statusC <- GAME_STATUS_RUNNING
			time.AfterFunc(bg.duringTime, func() {
				gr.statusC <- GAME_STATUS_FINISH
			})
		})
	})
}

type ClassRoomGame struct {
	scoreA         int
	scoreB         int
	lastHitA       time.Time
	lastHitB       time.Time
	lastPunishTime time.Time
	punishDuration time.Duration
	punishRole     int32
	BaseGame
}

func (crg *ClassRoomGame) OnEvent(gameRoomStatus int) {
	switch gameRoomStatus {
	case GAME_STATUS_READY:
		msg := GetGameMsg()
		msg.Event = game_event_ready
		msg.ReadyTime = 5
		msg.RoomNo = crg.gameRoom.index
		ants.Submit(func() {
			if msgA, errA := json.Marshal(msg); errA == nil {
				crg.gameRoom.BroadCast(msgA)
			}
		})
		ReturnGameMsg(msg)
	case GAME_STATUS_RUNNING:
		crg.startTime = time.Now()
		msg := GetGameMsg()
		msg.Event = game_event_start
		msg.ReadyTime = 0
		msg.Time = int(crg.duringTime / time.Second)
		msg.BScore = 0
		msg.AScore = 0
		if msgB, errB := json.Marshal(msg); errB == nil {
			crg.gameRoom.BroadCast(msgB)
		}
	case GAME_STATUS_FINISH:
		msg := GetGameMsg()
		msg.Event = game_event_finish
		msg.ReadyTime = 0
		msg.Time = int(time.Now().Sub(crg.startTime) / time.Second)
		msg.AScore = crg.scoreA
		msg.BScore = crg.scoreB
		if crg.scoreA > crg.scoreB {
			msg.Winner = "A"
		} else if crg.scoreA < crg.scoreB {
			msg.Winner = "B"
		} else {
			msg.Winner = "C"
		}
		if msgB, errB := json.Marshal(msg); errB == nil {
			crg.gameRoom.BroadCast(msgB)
		}

	}
}

func (crg *ClassRoomGame) Init(gr *GameRoom) {
	crg.gameRoom = gr
	crg.scoreA = 50
	crg.scoreB = 50
	crg.lastHitA = time.Now()
	crg.lastHitB = time.Now()
	crg.name = "教室战争"
	crg.code = "classroom war"
	crg.prefixTime = 3 * time.Second
	crg.duringTime = 50 * time.Second
	crg.punishDuration = 3 * time.Second
	crg.punishRole = ROOM_POS_EMPTY
	crg.npc = func(pos int32) interface{} {
		now := time.Now()
		rand.Seed(now.Unix())
		if now.Sub(crg.lastPunishTime) >= crg.punishDuration {
			if crg.punishRole != ROOM_POS_EMPTY {
				crg.punishRole = ROOM_POS_EMPTY
			}
			if ranNum := rand.Intn(100); ranNum%4 == 0 {
				//gr.BroadCast([]byte(EN1))
				crg.lastPunishTime = now
				if ranNum%2 == 0 {
					return game_event_npc2
				} else {
					return game_event_npc1
				}
			}
		}
		return ""
	}
	//crg.npc = func(pos int) bool{
	//	npcTimer := time.NewTicker(500 * time.Millisecond)
	//	for t := range npcTimer.C {
	//		rand.Seed(t.Unix())
	//		if time.Now().Sub(crg.lastPunishTime) >= 3*time.Second {
	//			if crg.punishRole != ROOM_POS_EMPTY {
	//				crg.punishRole = ROOM_POS_EMPTY
	//			}
	//			if rand.Intn(100)%4 == 0 {
	//				gr.BroadCast([]byte(EN1))
	//				crg.lastPunishTime = time.Now()
	//			}
	//		}
	//
	//	}
	//}
}
func (crg *ClassRoomGame) SetPunishRole(role int32) bool {
	return atomic.CompareAndSwapInt32(&crg.punishRole, crg.punishRole, role)
}
func (crg *ClassRoomGame) isPunish() bool {
	return time.Now().Sub(crg.lastPunishTime) < 3*time.Second

}
func (crg *ClassRoomGame) attack(client *GameClient, gr *GameRoom, gameMsg *GameMsg) *GameMsg {
	if crg.isPunish() {
		_, isMe := crg.checkPunish(client)
		if isMe {

		} else {
			if client.pos == ROOM_POS_A {
				crg.lastHitA = time.Now()
				gameMsg.AAttack = 1
				if crg.scoreB-2 <= 0 {
					crg.scoreB = 0
				} else {
					crg.scoreB -= 2
				}

			} else {
				crg.lastHitB = time.Now()
				gameMsg.BAttack = 1
				if crg.scoreA-2 <= 0 {
					crg.scoreA = 0
				} else {
					crg.scoreA -= 2
				}
			}
		}
	} else {
		npcEvent := crg.npc(client.pos)
		if npcEvent != "" {
			gameMsg.Event = npcEvent.(string)
			gameMsg.NpcObj = getPosStr(int(client.pos))
			crg.SetPunishRole(client.pos)
		} else {

		}
		if client.pos == ROOM_POS_A {
			crg.lastHitB = time.Now()
			gameMsg.AAttack = 1
			crg.scoreA += 1
		} else {
			crg.lastHitB = time.Now()
			gameMsg.BAttack = 1
			crg.scoreB = 1
		}

	}
	return gameMsg

}
func (crg *ClassRoomGame) checkPunish(client *GameClient) (int, bool) {
	if crg.punishRole == ROOM_POS_EMPTY {
		crg.SetPunishRole(client.pos)
		return -2, true
	} else if crg.punishRole == client.pos {
		return 0, true
	} else {
		return 2, false
	}
}

func (crg *ClassRoomGame) RunGame(gr *GameRoom) {
	crg.timeCounter(gr)
	ants.Submit(func() {
		ticker := time.NewTicker(1 * time.Second)
		for {
			select {
			case a := <-gr.clientA.read:
				gameMsg := GetGameMsg()
				if !gr.CheckStatus([]int32{GAME_STATUS_RUNNING}) {
					gameMsg.Code = 1

				} else {
					switch string(a) {
					case EA:
						if time.Now().Sub(crg.lastHitA) > 50*time.Millisecond {
							crg.lastHitA = time.Now()
							gameMsg.Event = game_event_attack
							crg.attack(gr.clientA, gr, gameMsg)
						}
					}
				}
				if gameMsg.Event != "" || gameMsg.Code != 0 {
					crg.fillScoreAndTime(gameMsg)
					b, err := json.Marshal(gameMsg)
					if err == nil {
						gr.BroadCast(b)
					}

				}

				ReturnGameMsg(gameMsg)
			case b := <-gr.clientB.read:
				gameMsg := GetGameMsg()
				if !gr.CheckStatus([]int32{GAME_STATUS_RUNNING}) {
					gr.BroadCast([]byte(GAME_ERROR_NOT_RUNNING))
				} else {
					switch string(b) {
					case EA:
						if time.Now().Sub(crg.lastHitB) > 50*time.Millisecond {
							crg.lastHitB = time.Now()
							gameMsg.Event = game_event_attack
							crg.attack(gr.clientB, gr, gameMsg)
						}

					}

				}
				if gameMsg.Event != "" || gameMsg.Code != 0 {
					crg.fillScoreAndTime(gameMsg)
					b, err := json.Marshal(gameMsg)
					if err == nil {
						gr.BroadCast(b)
					}
				}
				ReturnGameMsg(gameMsg)
			case <-ticker.C:
				if gr.GetStatus() == GAME_STATUS_FINISH { //停止监听
					break
				}

			}
		}
	})
}

func getPosStr(pos int) string {
	if pos == ROOM_POS_A {
		return "A"
	} else {
		return "B"
	}
}

func (crg *ClassRoomGame) fillScoreAndTime(gameMsg *GameMsg) {
	gameMsg.AScore = crg.scoreA
	gameMsg.BScore = crg.scoreB
	gameMsg.Time = int(time.Now().Sub(crg.startTime).Round(time.Second))
}
