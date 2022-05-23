package bot

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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

func (b *Bot) handleDeepLinkedStart(update *tgbotapi.Update) {
	payload := strings.TrimPrefix(update.Message.Text, "/start ")

	message := tgbotapi.NewMessage(update.Message.Chat.ID, "You have been invited to play a game!")
	message.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.InlineKeyboardButton{
				Text: "Join",
				WebApp: &tgbotapi.WebAppInfo{
					URL: fmt.Sprintf("https://127.0.0.1:8080/%s", payload),
				},
			}))

	b.api.Send(message)
}

func (b *Bot) handleStart(update *tgbotapi.Update) {
	message := "Welcome!\n" +
		"\n" +
		"Please use the main menu or add bot to a group chat."

	b.api.Send(tgbotapi.NewMessage(update.Message.Chat.ID, message))
}

func (b *Bot) handleCreate(update *tgbotapi.Update) {
	chat := update.Message.Chat
	if !chat.IsGroup() && !chat.IsSuperGroup() {
		message := "Please add bot to a group chat before creating a game."
		b.api.Send(tgbotapi.NewMessage(chat.ID, message))
		return
	}

	lobby, err := b.lobbyManager.CreateLobby()
	if err != nil {
		message := "Failed to create a game. Please try again."
		b.api.Send(tgbotapi.NewMessage(chat.ID, message))
		return
	}

	message := tgbotapi.NewMessage(chat.ID, "Successfully created a game!")
	deepLinkedUrl := fmt.Sprintf("https://t.me/%s?start=%x", b.api.Self.UserName, lobby.ID)
	message.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("Join game", deepLinkedUrl)))

	b.api.Send(message)
}

func (b *Bot) Run() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil || !update.Message.IsCommand() {
			continue
		}

		switch update.Message.Command() {
		case "start":
			if update.Message.Text != "/start" {
				b.handleDeepLinkedStart(&update)
			} else {
				b.handleStart(&update)
			}
		case "create":
			b.handleCreate(&update)
		default:
			b.api.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Unknown command"))
		}
	}

	return nil
}
