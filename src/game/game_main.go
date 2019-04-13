package game

import (
	"common/logging"
	"encoding/json"
	ws "github.com/gorilla/websocket"
	"github.com/panjf2000/ants"
	"net/http"
	"util"
	"util/uuid"
)

var (
	upgrader = ws.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		// 允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	GameMallInst = &GameMall{
		waitingClients: util.NewConcMap(),
		findClientId:   make(chan string, 1000),
		exitGameClient: make(chan *GameClient, 100),
	}
)

type GameMall struct {
	//游戏大厅，
	waitingClients util.ConcMap
	gameRooms      GameRooms
	findClientId   chan string
	exitGameClient chan *GameClient
}

func (gm *GameMall) addClient(client *GameClient) {
	logging.Info("add client id=%s", client.id)
	gm.waitingClients.Set(client.id, client)
	client.run()

}

func (gm *GameMall) Init() {
	GameMallInst.gameRooms = NewGameRooms()
	ants.Submit(func() {
		for {
			select {
			case clientId := <-gm.findClientId:
				logging.Info("get find req id=%s", clientId)
				if client, exist := gm.waitingClients.Pop(clientId); exist {
					if gr, err := gm.gameRooms.CreateOrGetGameRoom(); err == nil {
						gameClient := client.(*GameClient)
						gameClient.gameRoom = gr
						ants.Submit(gr.Run)
						gr.AddClient(gameClient)
						gameMsg := GetGameMsg()
						gameMsg.Event = game_event_match
						gameMsg.RoomNo = gr.index
						gameMsg.Pos = getPosStr(gameClient.pos)
						json, err := json.Marshal(gameMsg)
						if err == nil {
							gameClient.Send(json)
						}
						ReturnGameMsg(gameMsg)
					} else {
						client.(*GameClient).Send([]byte(GAME_ERROR_FIND))
					}
				}
			case client := <-gm.exitGameClient:
				gm.waitingClients.Set(client.id, client)
			}
		}
	})

}

//接入游戏,放个等候大厅
func GameDispatch(w http.ResponseWriter, r *http.Request) {
	//param := r.URL.Query()
	//logging.Debug(param.Get("appId"))

	if conn, err := upgrader.Upgrade(w, r, nil); err != nil {
		logging.Error("哦活,error:%s", err)
		conn.Close()
	} else {
		client := &GameClient{
			conn: conn,
			gm:   GameMallInst,
			send: make(chan []byte, 20),
			read: make(chan []byte),
			id:   uuid.Generate().String(),
		}
		GameMallInst.addClient(client)
	}

}
