package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/oevseev/gamebot/internal/lobby"
)

var bot *tgbotapi.BotAPI
var manager lobby.Manager

func handleStart(update *tgbotapi.Update) {
	bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Welcome!"))
}

func handleCreate(update *tgbotapi.Update) {
	lobby := manager.CreateLobby()
	log.Printf("Created lobby %x", lobby.ID)
	bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "OK!"))
}

func main() {
	token, ok := os.LookupEnv("TOKEN")
	if !ok {
		log.Panic("TOKEN not set")
	}

	var err error
	bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	manager = lobby.NewManager()

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
