package game

import (
	"common/logging"
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
		gameRooms:      func() GameRooms { return NewGameRooms() }(),
		findClientId:   make(chan string, 1000),
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
	gm.waitingClients.Set(client.id, client)
	client.gm = gm
	client.run()

}

func (gm *GameMall) Init() {
	ants.Submit(func() {
		for {
			select {
			case clientId := <-gm.findClientId:
				if client, exist := gm.waitingClients.Pop(clientId); exist {
					if gr, err := gm.gameRooms.CreateOrGetGameRoom(); err == nil {
						gr.Run()
						gr.AddClient(client.(*GameClient))
					} else {
						client.(*GameClient).Send([]byte(GAME_ERROR_FIND))
					}
				}



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
			send: make(chan []byte, 20),
			read: make(chan []byte),
			id:   uuid.Generate().String(),
		}
		GameMallInst.addClient(client)
	}

}


