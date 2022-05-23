package lobby

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/google/uuid"
)

type LobbyID uuid.UUID

type Lobby struct {
	ID      LobbyID
	Members map[int]tgbotapi.User
	manager *Manager
}

func (l *Lobby) Join(user tgbotapi.User) {
	l.Members[user.ID] = user
	l.manager.UpdateLobby(l)
}

func (l *Lobby) Leave(userId int) {
	delete(l.Members, userId)
	l.manager.UpdateLobby(l)
}
