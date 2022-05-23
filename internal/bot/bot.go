package bot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/oevseev/gamebot/internal/lobby"
	"go.mongodb.org/mongo-driver/mongo"
)

type Bot struct {
	api          *tgbotapi.BotAPI
	lobbyManager *lobby.Manager
}

func NewBot(api *tgbotapi.BotAPI, mongoClient *mongo.Client) *Bot {
	store := lobby.NewStore(mongoClient)
	manager := lobby.NewManager(store)

	return &Bot{
		api:          api,
		lobbyManager: manager,
	}
}

func (b *Bot) handleStart(update *tgbotapi.Update) {
	b.api.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Welcome!"))
}

func (b *Bot) handleCreate(update *tgbotapi.Update) {
	lobby, err := b.lobbyManager.CreateLobby()
	if err != nil {
		log.Printf("Failed to create lobby: %s", err)
		return
	}
	log.Printf("Created lobby %x", lobby.ID)

	b.api.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "OK!"))
}

func (b *Bot) Run() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := b.api.GetUpdatesChan(u)
	if err != nil {
		return err
	}

	for update := range updates {
		if update.Message == nil || !update.Message.IsCommand() {
			continue
		}

		switch update.Message.Command() {
		case "start":
			b.handleStart(&update)
		case "create":
			b.handleCreate(&update)
		default:
			b.api.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Unknown command"))
		}
	}

	return nil
}
