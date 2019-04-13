package game

import (
	"common/logging"
	"encoding/json"
	"github.com/panjf2000/ants"
	"math/rand"
	"sync/atomic"
	"time"
	"util"
)

const (
	//EA  = "EA"  //攻击事件
	EA  = "attack" //攻击事件
	EAA = "EAA"    //A攻击事件
	EAB = "EAB"    //B攻击事件
	EN1 = "EN1"    //npc事件1
	EN2 = "EN2"    //npc事件2
	OB1 = "OB1"    //掉血1
	OB2 = "OB2"    //掉血2
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
	npc        func(param interface{}) interface{}
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
			if gr.CheckStatus([]int32{GAME_STATUS_READY}) {
				gr.statusC <- GAME_STATUS_RUNNING
			}
			time.AfterFunc(bg.duringTime, func() {
				if gr.CheckStatus([]int32{GAME_STATUS_RUNNING}) {
					gr.statusC <- GAME_STATUS_FINISH
				}
			})
		})
	})
}

type ClassRoomGame struct {
	scoreA         int
	scoreB         int
	lastHitA       time.Time
	lastHitB       time.Time
	lastCheckTime  time.Time
	lastPunishTime time.Time
	punishDuration time.Duration
	punishRole     int32
	BaseGame
}

func (crg *ClassRoomGame) OnEvent(gameRoomStatus int) {
	switch gameRoomStatus {
	case GAME_STATUS_READY:
		msgA, msgB := GetGameMsg(), GetGameMsg()
		msgA.Event, msgB.Event = game_event_ready, game_event_ready
		msgA.ReadyTime, msgB.ReadyTime = 5, 5
		msgA.RoomNo, msgB.RoomNo = crg.gameRoom.index, crg.gameRoom.index
		msgA.OppName = crg.gameRoom.clientB.name
		msgB.OppName = crg.gameRoom.clientA.name
		crg.SendAAndBAndReturn(msgA, msgB, crg.gameRoom)
	case GAME_STATUS_RUNNING:
		crg.startTime = time.Now()
		crg.npc(nil)
		msg := GetGameMsg()
		msg.Event = game_event_start
		msg.ReadyTime = 0
		msg.Time = int(crg.duringTime / time.Second)
		msg.BScore = 0
		msg.AScore = 0
		crg.BroadcastAndReturn(msg, crg.gameRoom)
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
			msg.Winner = "C" //平局
		}
		crg.BroadcastAndReturn(msg, crg.gameRoom)
	case GAME_STATUS_FINISH_ERROR:
		msg := GetGameMsg()
		msg.Event = game_event_finish_error
		crg.BroadcastAndReturn(msg, crg.gameRoom)
	}
}

func (crg *ClassRoomGame) Init(gr *GameRoom) {
	now := time.Now()
	crg.gameRoom = gr
	crg.scoreA = 0
	crg.scoreB = 0
	crg.lastHitA = now
	crg.lastHitB = now
	crg.lastCheckTime = now
	crg.name = "教室战争"
	crg.code = "classroom war"
	crg.prefixTime = 3 * time.Second
	crg.duringTime = 60 * time.Second
	crg.punishDuration = 2 * time.Second
	crg.punishRole = ROOM_POS_EMPTY
	crg.npc = func(param interface{}) interface{} {
		ants.Submit(func() {
			npcTimer := time.NewTicker(500 * time.Millisecond)
			count := 0
			for now := range npcTimer.C {
				punishMinus := util.Abs(int64(now.Sub(crg.lastPunishTime).Round(time.Millisecond) - crg.punishDuration))
				if punishMinus < int64(100*time.Millisecond) {
					gameMsg := GetGameMsg()
					gameMsg.Event = game_event_npc3
					crg.fillAndBroadcastAndReturn(gameMsg, gr)
				} else if now.Sub(crg.lastPunishTime) >= crg.punishDuration && crg.getTimeLeft(time.Now()) > 5 { //最后5秒不产生惩罚
					gameMsg := GetGameMsg()
					rand.Seed(now.Unix())
					randNum := rand.Intn(100)
					checkMinus := util.Abs(int64(now.Sub(crg.lastCheckTime).Round(time.Millisecond) - 3*time.Second))
					//logging.Info("the minus is %d", checkMinus)
					if checkMinus < int64(100*time.Millisecond) {
						crg.lastPunishTime = now
						gameMsg.Event = util.AOrB(func() bool {
							return randNum%2 == 0
						}, game_event_npc2, game_event_npc22).(string)
						if now.Sub(crg.lastHitA) < 500*time.Millisecond || now.Sub(crg.lastHitB) < 500*time.Millisecond {
							if now.Sub(crg.lastHitA) < 500*time.Millisecond && now.Sub(crg.lastHitB) < 500*time.Millisecond {
								gameMsg.NpcObj = getPosStr(ROOM_POS_ALL)
								crg.SetPunishRole(ROOM_POS_ALL)
							} else if now.Sub(crg.lastHitA) < 500*time.Millisecond {
								gameMsg.NpcObj = getPosStr(ROOM_POS_A)
								crg.SetPunishRole(ROOM_POS_A)
							} else if now.Sub(crg.lastHitB) < 500*time.Millisecond {
								gameMsg.NpcObj = getPosStr(ROOM_POS_B)
								crg.SetPunishRole(ROOM_POS_B)
							}
						} else {
							gameMsg.NpcObj = getPosStr(ROOM_POS_EMPTY)
							crg.SetPunishRole(ROOM_POS_EMPTY)
						}
						crg.fillAndBroadcastAndReturn(gameMsg, gr)
					} else {
						if ((randNum/2 == 0 && now.Sub(crg.lastPunishTime) > 2*time.Second) || now.Sub(crg.lastPunishTime) >= 6*time.Second) && count < 9 && now.Sub(crg.lastCheckTime) > 3*time.Second && now.Sub(crg.startTime) > 3*time.Second {
							count++
							crg.lastCheckTime = now
							gameMsg := GetGameMsg()
							gameMsg.Event = game_event_npc1
							crg.fillAndBroadcastAndReturn(gameMsg, gr)
						} else {
							if crg.punishRole != ROOM_POS_EMPTY {
								crg.punishRole = ROOM_POS_EMPTY
							}
						}

					}

				} else {
					if crg.gameRoom.CheckStatus([]int32{GAME_STATUS_FINISH, GAME_STATUS_FINISH_ERROR}) {
						return
					}
				}
			}
		})

		//if now.Sub(crg.lastPunishTime) >= crg.punishDuration {
		//	if crg.punishRole != ROOM_POS_EMPTY {
		//		crg.punishRole = ROOM_POS_EMPTY
		//	}
		//	if ranNum := rand.Intn(100); ranNum%4 == 0 {
		//		//gr.BroadCast([]byte(EN1))
		//		crg.lastPunishTime = now
		//		if ranNum%2 == 0 {
		//			return game_event_npc2
		//		} else {
		//			return game_event_npc1
		//		}
		//	}
		//}
		return nil
	}
}
func (crg *ClassRoomGame) SetPunishRole(role int32) bool {
	return atomic.CompareAndSwapInt32(&crg.punishRole, crg.punishRole, role)
}
func (crg *ClassRoomGame) isPunish() bool {
	return time.Now().Sub(crg.lastPunishTime) < crg.punishDuration

}
func (crg *ClassRoomGame) attack(client *GameClient, gr *GameRoom, gameMsg *GameMsg) *GameMsg {
	if crg.isPunish() {
		step, isMe := crg.checkPunish(client)
		if isMe {

		} else {
			gameMsg.Event = game_event_attack
			if client.pos == ROOM_POS_A {
				crg.lastHitA = time.Now()
				gameMsg.AAttack = 1
				gameMsg.BAttack = step
				if crg.scoreB+step <= 0 {
					crg.scoreB = 0
				} else {
					crg.scoreB += step
				}
			} else {
				crg.lastHitB = time.Now()
				gameMsg.BAttack = 1
				gameMsg.AAttack = step
				if crg.scoreA+step <= 0 {
					crg.scoreA = 0
				} else {
					crg.scoreA += step
				}
			}
		}
	} else {
		gameMsg.Event = game_event_attack
		if client.pos == ROOM_POS_A {
			crg.lastHitA = time.Now()
			gameMsg.AAttack = 1
			crg.scoreA += 1
		} else {
			crg.lastHitB = time.Now()
			gameMsg.BAttack = 1
			crg.scoreB += 1
		}

	}
	return gameMsg

}
func (crg *ClassRoomGame) checkPunish(client *GameClient) (int, bool) {
	if crg.punishRole == client.pos || crg.punishRole == ROOM_POS_ALL {
		return 0, true
	} else {
		return -2, false
	}
}
func (crg *ClassRoomGame) Broadcast(msg *GameMsg, gr *GameRoom) {
	if msg.Event != "" || msg.Code != 0 {
		b, err := json.Marshal(msg)
		if err == nil {
			gr.BroadCast(b)
		}
	}
}
func (crg *ClassRoomGame) BroadcastAndReturn(msg *GameMsg, gr *GameRoom) {
	crg.Broadcast(msg, gr)
	ReturnGameMsg(msg)
}
func (crg *ClassRoomGame) SendAAndBAndReturn(msgA *GameMsg, msgB *GameMsg, gr *GameRoom) {
	a, b := []byte(""), []byte("")
	if msgA.Event != "" || msgA.Code != 0 {
		a, _ = json.Marshal(msgA)
	}
	if msgB.Event != "" || msgB.Code != 0 {
		b, _ = json.Marshal(msgB)
	}
	ants.Submit(func() {
		gr.clientA.Send(a)
	})
	ants.Submit(func() {
		gr.clientB.Send(b)
	})
	ReturnGameMsg(msgA, msgB)
}

func (crg *ClassRoomGame) fillAndBroadcastAndReturn(gameMsg *GameMsg, gr *GameRoom) {
	crg.fillScoreAndTime(gameMsg)
	crg.Broadcast(gameMsg, gr)
	ReturnGameMsg(gameMsg)
}

func (crg *ClassRoomGame) RunGame(gr *GameRoom) {
	crg.timeCounter(gr)
	ants.Submit(func() {
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

							crg.attack(gr.clientA, gr, gameMsg)
						}
					}
				}
				crg.fillAndBroadcastAndReturn(gameMsg, gr)
			case b := <-gr.clientB.read:
				gameMsg := GetGameMsg()
				if !gr.CheckStatus([]int32{GAME_STATUS_RUNNING}) {
					gameMsg.Code = 1
				} else {
					switch string(b) {
					case EA:
						if time.Now().Sub(crg.lastHitB) > 50*time.Millisecond {
							gameMsg.Event = game_event_attack
							crg.attack(gr.clientB, gr, gameMsg)
						}
					}

				}
				crg.fillAndBroadcastAndReturn(gameMsg, gr)
			}
		}
	})
}

func getPosStr(pos int32) string {
	switch pos {
	case ROOM_POS_A:
		return "A"
	case ROOM_POS_B:
		return "B"
	case ROOM_POS_ALL:
		return "C"
	case ROOM_POS_EMPTY:
		return "D"

	}
	return "D"
}

func (crg *ClassRoomGame) fillScoreAndTime(gameMsg *GameMsg) {
	gameMsg.AScore = crg.scoreA
	gameMsg.BScore = crg.scoreB
	gameMsg.Time = crg.getTimeLeft(time.Now())
}
func (crg *ClassRoomGame) getTimeLeft(input time.Time) int {
	min := crg.duringTime - input.Sub(crg.startTime)
	return int(min.Round(time.Second) / time.Second)
}
