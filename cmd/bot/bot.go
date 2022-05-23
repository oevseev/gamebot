package main

import (
	"context"
	"log"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/oevseev/gamebot/internal/lobby"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var bot *tgbotapi.BotAPI
var manager *lobby.Manager

func handleStart(update *tgbotapi.Update) {
	bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Welcome!"))
}

func handleCreate(update *tgbotapi.Update) {
	lobby, err := manager.CreateLobby()
	if err != nil {
		log.Printf("Failed to create lobby: %s", err)
		return
	}
	log.Printf("Created lobby %x", lobby.ID)
	bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "OK!"))
}

func main() {
	token, ok := os.LookupEnv("BOT_API_TOKEN")
	if !ok {
		log.Panic("TOKEN not set")
	}

	mongoEndpoint, ok := os.LookupEnv("MONGO_ENDPOINT")
	if !ok {
		log.Panic("MONGO_ENDPOINT not set")
	}

	var err error
	bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoEndpoint))
	if err != nil {
		log.Panic(err)
	}

	store := lobby.NewStore(client)
	manager = lobby.NewManager(store)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if !update.Message.IsCommand() {
			continue
		}

		switch update.Message.Command() {
		case "start":
			handleStart(&update)
		case "create":
			handleCreate(&update)
		default:
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Unknown command"))
		}
	}
}
