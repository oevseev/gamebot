package lobby

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/google/uuid"
)

type LobbyID uuid.UUID

type Lobby struct {
	ID      LobbyID
	Members map[int]tgbotapi.User
}

func NewLobby() *Lobby {
	return &Lobby{
		ID:      LobbyID(uuid.New()),
		Members: make(map[int]tgbotapi.User),
	}
}

func (l *Lobby) Join(user tgbotapi.User) {
	l.Members[user.ID] = user
}

func (l *Lobby) Leave(userId int) {
	delete(l.Members, userId)
}
