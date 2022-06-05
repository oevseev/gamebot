package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/oevseev/gamebot/internal/webserver"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	mongoEndpoint, ok := os.LookupEnv("MONGO_ENDPOINT")
	if !ok {
		log.Panic("MONGO_ENDPOINT not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoEndpoint))
	if err != nil {
		log.Panic(err)
	}

	w := webserver.NewWebServer(fqdn, mongoClient)
	if err := w.RunTLS(listenAddr, tlsCertPath, tlsKeyPath); err != nil {
		log.Panic(err)
	}
}
