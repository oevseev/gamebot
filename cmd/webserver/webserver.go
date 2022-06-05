package main

import (
	"log"
	"os"

	"github.com/oevseev/gamebot/internal/webserver"
)

func main() {
	fqdn, ok := os.LookupEnv("FQDN")
	if !ok {
		log.Panic("FQDN not set")
	}

	listenAddr, ok := os.LookupEnv("LISTEN_ADDR")
	if !ok {
		log.Panic("LISTEN_ADDR not set")
	}

	tlsCertPath, ok := os.LookupEnv("TLS_CERT_PATH")
	if !ok {
		log.Panic("TLS_CERT_PATH not set")
	}

	tlsKeyPath, ok := os.LookupEnv("TLS_KEY_PATH")
	if !ok {
		log.Panic("TLS_KEY_PATH not set")
	}

	w := webserver.NewWebServer(fqdn)
	if err := w.RunTLS(listenAddr, tlsCertPath, tlsKeyPath); err != nil {
		log.Panic(err)
	}
}
