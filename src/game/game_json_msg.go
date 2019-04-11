package game

import "sync"

const (
	game_event_match  = "match"
	game_event_ready  = "ready"
	game_event_start  = "start"
	game_event_attack = "attack"
	game_event_npc1   = "npc1"
	game_event_npc2   = "npc2"
	game_event_finish = "finish"
)

type GameMsg struct {
	Code      int    `json:"code"`                 //0为正常，1为错误
	Event     string `json:"event"`                //分别为ready,attack,npc,finish
	ReadyTime int    `json:"readyTime,omitempty"`  //准备时间(秒)
	Time      int    `json:"time"`                 //剩余时间(秒)
	AAttack   int    `json:"AAttack"`              //A角色攻击
	BAttack   int    `json:"BAttack"`              //A角色攻击
	AScore    int    `json:"AScore"`               //A角色得分
	BScore    int    `json:"BScore"`               //A角色得分
	NpcObj    string `json:"npcObj,omitempty"`     //惩罚角色，分别A,B
	Pos       string  `json:"pos,string,omitempty"` //座位位置
	RoomNo    int    `json:"roomNo,omitempty"`     //房间号
	Winner    string `json:"winner,omitempty"`     //胜者
}

//一个pp的对象池
var GameMsgPool = sync.Pool{
	New: func() interface{} { return new(GameMsg) },
}
// 分配一个新的pp或者拿一个缓存的。
func GetGameMsg() *GameMsg {
	m := GameMsgPool.Get().(*GameMsg)
	m.Event = ""
	m.ReadyTime = 0
	m.Time = 0
	m.AAttack = 0
	m.BAttack = 0
	m.AScore = 0
	m.BScore = 0
	m.NpcObj = ""
	m.Pos = ""
	m.RoomNo = 0
	m.Winner = ""
	return m
}
func ReturnGameMsg(msgs ...*GameMsg) {
	for _, v := range msgs {
		GameMsgPool.Put(v)
	}

}
