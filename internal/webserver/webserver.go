package webserver

import (
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/oevseev/gamebot/internal/lobby"
	"go.mongodb.org/mongo-driver/mongo"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebServer struct {
	r            *gin.Engine
	lobbyManager *lobby.Manager
	gameServer   *GameServer
}

func NewWebServer(fqdn string, mongoClient *mongo.Client) *WebServer {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*.tmpl")

	store := lobby.NewStore(mongoClient)
	manager := lobby.NewManager(store)

	w := &WebServer{
		r:            r,
		lobbyManager: manager,
		gameServer:   NewGameServer(manager),
	}

	r.GET("/ws", func(c *gin.Context) {
		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer ws.Close()
		w.gameServer.HandleConnection(ws)
	})

	r.GET("/:id", func(c *gin.Context) {
		hex, err := hex.DecodeString(c.Param("id"))
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		uuid, err := uuid.FromBytes(hex)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		if _, err := w.lobbyManager.GetLobby(lobby.ID(uuid)); err != nil {
			c.AbortWithError(http.StatusNotFound, err)
			return
		}
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"webSocketUrl": fmt.Sprintf("wss://%s/ws", fqdn),
			"gameId":       c.Param("id"),
		})
	})

	return w
}

func (w *WebServer) RunTLS(listenAddr string, tlsCertPath string, tlsKeyPath string) error {
	return w.r.RunTLS(listenAddr, tlsCertPath, tlsKeyPath)
}
