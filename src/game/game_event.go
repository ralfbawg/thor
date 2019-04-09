package game

import (
	"common/logging"
	"github.com/panjf2000/ants"
	"math/rand"
	"strconv"
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
}

type BaseGame struct {
	name       string
	code       string
	prefixTime time.Duration
	duringTime time.Duration
	npc        func()
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
			ants.Submit(bg.npc)
		})
	})
}

type ClassRoomGame struct {
	bloodA         int
	bloodB         int
	lastHitA       time.Time
	lastHitB       time.Time
	lastPunishTime time.Time
	punishDuration time.Duration
	punishRole     int32
	*BaseGame
}

func (crg *ClassRoomGame) Init(gr *GameRoom) {
	crg.bloodA = 30
	crg.bloodB = 30
	crg.lastHitA = time.Now()
	crg.lastHitB = time.Now()
	crg.name = "教室战争"
	crg.code = "classroom war"
	crg.prefixTime = 3 * time.Second
	crg.duringTime = 30 * time.Second
	crg.punishDuration = 3 * time.Second
	crg.punishRole = ROOM_POS_EMPTY
	crg.npc = func() {
		npcTimer := time.NewTicker(500 * time.Millisecond)
		for t := range npcTimer.C {
			rand.Seed(t.Unix())
			if time.Now().Sub(crg.lastPunishTime) >= 3*time.Second {
				if crg.punishRole != ROOM_POS_EMPTY {
					crg.punishRole = ROOM_POS_EMPTY
				}
				if rand.Intn(100)%4 == 0 {
					gr.BroadCast([]byte(EN1))
					crg.lastPunishTime = time.Now()
				}
			}

		}
	}
}
func (crg *ClassRoomGame) SetPunishRole(role int32) bool {
	return atomic.CompareAndSwapInt32(&crg.punishRole, crg.punishRole, role)
}
func (crg *ClassRoomGame) isPunish() bool {
	return time.Now().Sub(crg.lastPunishTime) < 3*time.Second

}
func (crg *ClassRoomGame) attack(client *GameClient, result string, gr *GameRoom) string {
	if crg.isPunish() {
		crg.checkPunish(client)
	} else {
		crg.lastHitA = time.Now()
		result = EAA + "," + strconv.Itoa(crg.bloodA) + "," + strconv.Itoa(crg.bloodB)
		if RemoveBloodAndCheckFinish(&crg.bloodB, gr) {
			result += "," + GAME_EVENT_FINISH + "," + GAME_EVNET_WINNER_A
		}
	}
	return result

}
func (crg *ClassRoomGame) checkPunish(client *GameClient) (int, bool, bool) {
	if crg.punishRole == ROOM_POS_EMPTY {
		crg.SetPunishRole(client.pos)
		return -2, true, true
	} else if crg.punishRole == client.pos {
		return 0, true, true
	} else {
		return 2, false, true
	}
}

func (crg *ClassRoomGame) RunGame(gr *GameRoom) {
	crg.timeCounter(gr)
	ants.Submit(func() {
		ticker := time.NewTicker(1 * time.Second)
		for {
			select {
			case a := <-gr.clientA.read:
				if !gr.CheckStatus([]int32{GAME_STATUS_RUNNING}) {
					gr.BroadCast([]byte(GAME_ERROR_NOT_RUNNING))
				} else {
					if time.Now().Sub(crg.lastHitA) > 50*time.Millisecond {
						crg.lastHitA = time.Now()
						back := ""
						switch string(a) {
						case EA:
							back := crg.attack(gr.clientA, back, gr)

						}
						gr.BroadCast([]byte(back))

						if crg.isPunish() {

						} else {
							crg.lastHitA = time.Now()
							result = EAA + "," + strconv.Itoa(crg.bloodA) + "," + strconv.Itoa(crg.bloodB)
							if RemoveBloodAndCheckFinish(&crg.bloodB, gr) {
								result += "," + GAME_EVENT_FINISH + "," + GAME_EVNET_WINNER_A
							}
						}

						blood, isMe, valid := crg.checkPunish(ROOM_POS_A)
						if valid {
							if isMe {
								if RemoveBloodAndCheckFinish(&crg.bloodA, gr) {
									result += "," + GAME_EVENT_FINISH + "," + GAME_EVNET_WINNER_B
								}
							} else {
								if blood != 0 && RemoveBloodAndCheckFinish(&crg.bloodB, gr) {
									result += "," + GAME_EVENT_FINISH + "," + GAME_EVNET_WINNER_A
								}
							}
						}

						if time.Now().Sub(crg.lastPunishTime) < crg.punishDuration {
							if crg.punishRole == -1 {
								crg.SetPunishRole(ROOM_POS_A)
							} else if crg.punishRole == ROOM_POS_A {
								result = EAA + "," + strconv.Itoa(crg.bloodA) + "," + strconv.Itoa(crg.bloodB)

							} else {

							}

						} else {
							crg.lastHitA = time.Now()
							result = EAA + "," + strconv.Itoa(crg.bloodA) + "," + strconv.Itoa(crg.bloodB)
							if RemoveBloodAndCheckFinish(&crg.bloodB, gr) {
								result += "," + GAME_EVENT_FINISH + "," + GAME_EVNET_WINNER_A
							}
						}
						gr.BroadCast([]byte(result))
					}
				}
			case b := <-gr.clientB.read:
				if !gr.CheckStatus([]int32{GAME_STATUS_RUNNING}) {
					gr.BroadCast([]byte(GAME_ERROR_NOT_RUNNING))
				} else {
					if string(b) == EA && time.Now().Sub(crg.lastHitB) > 50*time.Millisecond {
						if crg.punishRole == -1 {
							crg.SetPunishRole(ROOM_POS_B)
						}
						result := ""
						if time.Now().Sub(crg.lastPunishTime) < crg.punishDuration { //惩罚
							result = EAB + "," + strconv.Itoa(crg.bloodA) + "," + strconv.Itoa(crg.bloodB)
							if RemoveBloodAndCheckFinish(&crg.bloodB, gr) {
								result += "," + GAME_EVENT_FINISH + "," + GAME_EVNET_WINNER_A
							}
						} else {
							crg.lastHitB = time.Now()
							result = EAB + "," + strconv.Itoa(crg.bloodA) + "," + strconv.Itoa(crg.bloodB)
							if RemoveBloodAndCheckFinish(&crg.bloodA, gr) {
								result += "," + GAME_EVENT_FINISH + "," + GAME_EVNET_WINNER_B
							}

						}
						gr.BroadCast([]byte(result))
					}
				}
			case <-ticker.C:
				if gr.GetStatus() == GAME_STATUS_FINISH { //停止监听
					break
				}

			}
		}
	})
}
func RemoveBloodAndCheckFinish(blood *int, gr *GameRoom) bool { //只有一个入口扣减，应该不需要同步
	*blood -= 2
	if *blood <= 0 {
		gr.BroadCast([]byte(GAME_EVENT_FINISH))
		gr.statusC <- GAME_STATUS_FINISH
		return true
	} else {
		return false
	}
}
