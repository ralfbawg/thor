package game

import (
	"common/logging"
	"github.com/panjf2000/ants"
	"math/rand"
	"strconv"
	"time"
)

const (
	EA  = "EA"  //攻击
	EAA = "EAA" //A攻击
	EAB = "EAB" //B攻击
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
	crg.npc = func() {
		npcTimer := time.NewTicker(500 * time.Millisecond)
		for t := range npcTimer.C {
			rand.Seed(t.Unix())
			if rand.Intn(100)%4 == 0 {
				gr.BroadCast([]byte(EN1))
				crg.lastPunishTime = time.Now()
			}
		}
	}
}
func (crg *ClassRoomGame) RunGame(gr *GameRoom) {
	crg.timeCounter(gr)
	ants.Submit(func() {
		for {
			select {
			case a := <-gr.clientA.read:
				if !gr.CheckStatus([]int32{GAME_STATUS_RUNNING}) {
					gr.BroadCast([]byte(GAME_ERROR_NOT_RUNNING))
				} else {
					if string(a) == EA && time.Now().Sub(crg.lastHitA) > 50*time.Millisecond {
						result := ""
						if time.Now().Sub(crg.lastPunishTime) < crg.punishDuration {
							result = EAA + "," + strconv.Itoa(crg.bloodA) + "," + strconv.Itoa(crg.bloodB)
							if RemoveBloodAndCheckFinish(&crg.bloodA, gr) {
								result += "," + GAME_EVENT_FINISH + "," + GAME_EVNET_WINNER_B
							}
						} else {
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
					if string(b) == EA && time.Now().Sub(crg.lastHitA) > 50*time.Millisecond {
						result := ""
						if time.Now().Sub(crg.lastPunishTime) < crg.punishDuration { //惩罚
							result = EAB + "," + strconv.Itoa(crg.bloodA) + "," + strconv.Itoa(crg.bloodB)
							if RemoveBloodAndCheckFinish(&crg.bloodB, gr) {
								result += "," + GAME_EVENT_FINISH + "," + GAME_EVNET_WINNER_A
							}
						} else {
							result = EAB + "," + strconv.Itoa(crg.bloodA) + "," + strconv.Itoa(crg.bloodB)
							if RemoveBloodAndCheckFinish(&crg.bloodA, gr) {
								result += "," + GAME_EVENT_FINISH + "," + GAME_EVNET_WINNER_B
							}

						}
						gr.BroadCast([]byte(result))
					}
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
