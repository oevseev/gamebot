package webserver

import (
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/oevseev/gamebot/internal/games/card/preferans"
	"github.com/oevseev/gamebot/internal/lobby"
)

type GameServer struct {
	lobbyManager *lobby.Manager
	games        map[lobby.ID]*preferans.Game
}

func NewGameServer(lobbyManager *lobby.Manager) *GameServer {
	return &GameServer{
		lobbyManager: lobbyManager,
		games:        make(map[lobby.ID]*preferans.Game),
	}
}

func (g *GameServer) HandleConnection(ws *websocket.Conn) {
	for {
		mt, message, err := ws.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}
		if mt != websocket.TextMessage {
			continue
		}
		if err := g.handleMessage(ws, message); err != nil {
			log.Println(err)
		}
	}
}

type Message struct {
	MessageType string                 `json:"messageType"`
	Payload     map[string]interface{} `json:"payload"`
}

func (g *GameServer) handleMessage(ws *websocket.Conn, messageRaw []byte) error {
	var message Message
	err := json.Unmarshal(messageRaw, &message)
	if err != nil {
		return err
	}

	switch message.MessageType {
	case "authorize":
		gameId, ok := message.Payload["gameId"].(string)
		if !ok {
			return errors.New("gameId is not string")
		}
		lobbyId, err := lobby.IDFromString(gameId)
		if err != nil {
			return err
		}
		playerIdString, ok := message.Payload["playerId"].(string)
		if !ok {
			return errors.New("playerId is not string")
		}
		playerId, err := strconv.ParseInt(playerIdString, 10, 64)
		if err != nil {
			return err
		}
		if err := g.handleAuthorize(ws, lobbyId, playerId); err != nil {
			return err
		}
	}

	return nil
}

// TODO: Devise a better solution when store is thread safe
var autorizeMutex sync.Mutex

func (g *GameServer) handleAuthorize(ws *websocket.Conn, gameId lobby.ID, playerId int64) error {
	autorizeMutex.Lock()

	lobby, err := g.lobbyManager.GetLobby(gameId)
	if err != nil {
		return err
	}

	_, isMember := lobby.Members[playerId]
	if !isMember && len(lobby.Members) < 3 {
		if err := lobby.Join(playerId); err != nil {
			return err
		}
		isMember = true
	}

	autorizeMutex.Unlock()

	var message Message
	if isMember {
		message = Message{
			MessageType: "joinedAsPlayer",
			Payload: map[string]interface{}{
				"preferansConfig": g.games[gameId].Config,
				"preferansState":  g.games[gameId].GetPublicState(playerId),
			},
		}
	} else {
		message = Message{
			MessageType: "joinedAsSpectator",
			Payload:     map[string]interface{}{},
		}
	}
	payload, err := json.Marshal(message)
	if err != nil {
		return err
	}
	ws.WriteMessage(websocket.TextMessage, payload)

	return nil
}
