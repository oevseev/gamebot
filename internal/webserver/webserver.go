package webserver

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
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

func ShowErrorPage(c *gin.Context) {
	c.Next()

	if len(c.Errors) > 0 {
		c.HTML(-1, "error.tmpl", nil)
	}
}

func NewWebServer(fqdn string, mongoClient *mongo.Client) *WebServer {
	r := gin.Default()
	r.Use(ShowErrorPage)
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
			log.Println(err)
			return
		}
		defer ws.Close()
		w.gameServer.HandleConnection(ws)
	})

	r.GET("/:id", func(c *gin.Context) {
		lobbyId, err := lobby.IDFromString(c.Param("id"))
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		if _, err := w.lobbyManager.GetLobby(lobbyId); err != nil {
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
