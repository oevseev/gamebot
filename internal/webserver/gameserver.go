package webserver

import (
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/oevseev/gamebot/internal/games/card/preferans"
	"github.com/oevseev/gamebot/internal/lobby"
)

type GameServer struct {
	lobbyManager *lobby.Manager
	games        map[string]preferans.Game
}

func NewGameServer(lobbyManager *lobby.Manager) *GameServer {
	return &GameServer{
		lobbyManager: lobbyManager,
		games:        map[string]preferans.Game{},
	}
}

func (g *GameServer) HandleConnection(ws *websocket.Conn) {
	for {
		mt, message, err := ws.ReadMessage()
		if err != nil {
			fmt.Println(err)
			break
		}
		if mt != websocket.TextMessage {
			continue
		}
		fmt.Println(string(message))
	}
}
