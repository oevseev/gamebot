package lobby

import (
	"github.com/google/uuid"
)

type ID uuid.UUID

type Lobby struct {
	ID      ID
	Members map[int]struct{}
	manager *Manager
}

func (l *Lobby) Join(userId int) {
	l.Members[userId] = struct{}{}
	l.manager.updateLobby(l)
}

func (l *Lobby) Leave(userId int) {
	delete(l.Members, userId)
	l.manager.updateLobby(l)
}
