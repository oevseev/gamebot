package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.String(http.StatusOK, id)
	})

	r.RunTLS(
		os.Getenv("LISTEN_ADDR"),
		os.Getenv("TLS_CERT_PATH"),
		os.Getenv("TLS_KEY_PATH"),
	)
}
