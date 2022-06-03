package main

import (
	"context"
	"log"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/oevseev/gamebot/internal/bot"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	token, ok := os.LookupEnv("BOT_API_TOKEN")
	if !ok {
		log.Panic("TOKEN not set")
	}

	mongoEndpoint, ok := os.LookupEnv("MONGO_ENDPOINT")
	if !ok {
		log.Panic("MONGO_ENDPOINT not set")
	}

	webappUrl, ok := os.LookupEnv("WEBAPP_URL")
	if !ok {
		log.Panic("WEBAPP_URL not set")
	}

	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}
	api.Debug = true
	log.Printf("Authorized on account %s", api.Self.UserName)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoEndpoint))
	if err != nil {
		log.Panic(err)
	}

	b := bot.NewBot(api, mongoClient, webappUrl)
	if err := b.Run(); err != nil {
		log.Panic(err)
	}
}
