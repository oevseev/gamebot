package webserver

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebServer struct {
	r *gin.Engine
}

func NewWebServer(fqdn string) *WebServer {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*.tmpl")

	r.GET("/ws", func(c *gin.Context) {
		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer ws.Close()
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
	})

	r.GET("/:id", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"webSocketUrl": fmt.Sprintf("wss://%s/ws", fqdn),
			"gameId":       c.Param("id"),
		})
	})

	return &WebServer{r: r}
}

func (w *WebServer) RunTLS(listenAddr string, tlsCertPath string, tlsKeyPath string) error {
	return w.r.RunTLS(listenAddr, tlsCertPath, tlsKeyPath)
}
