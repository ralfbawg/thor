package game

const (
	game_event_ready  = "ready"
	game_event_start  = "start"
	game_event_attack = "attack"
	game_event_npc    = "npc"
	game_event_finish = "finish"
)

type GameMsg struct {
	Event     string `json:"event"`            //分别为ready,attack,npc,finish
	ReadyTime int    `json:"readyTime"`        //准备时间(秒)
	Time      int    `json:"time"`             //剩余时间(秒)
	AAttack   int    `json:"AAttack"`           //A角色得分
	BAttack   int    `json:"BAttack"`           //A角色得分
	AScore    int    `json:"AScore"`           //A角色得分
	BScore    int    `json:"BScore"`           //A角色得分
	NpcObj    string `json:"npcObj,omitempty"` //惩罚角色，分别A,B
	Pos       string `json:"pos,omitempty"`    //座位位置
	RoomNo    int    `json:"roomNo,omitempty"` //房间号
	Winner    string `json:"winner,omitempty"` //胜者
}
