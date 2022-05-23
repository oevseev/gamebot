package main

import (
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.RunTLS(
		os.Getenv("LISTEN_ADDR"),
		os.Getenv("TLS_CERT_PATH"),
		os.Getenv("TLS_KEY_PATH"),
	)
}
