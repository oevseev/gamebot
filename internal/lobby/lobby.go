package lobby

import (
	"encoding/hex"

	"github.com/google/uuid"
)

type ID uuid.UUID

func IDFromString(s string) (ID, error) {
	hex, err := hex.DecodeString(s)
	if err != nil {
		return ID{}, err
	}
	id, err := uuid.FromBytes(hex)
	if err != nil {
		return ID{}, err
	}
	return ID(id), nil
}

type Lobby struct {
	ID      ID
	Members map[int64]struct{}
	manager *Manager
}

func (l *Lobby) Join(userId int64) error {
	l.Members[userId] = struct{}{}
	if err := l.manager.updateLobby(l); err != nil {
		return err
	}
	return nil
}

func (l *Lobby) Leave(userId int64) error {
	delete(l.Members, userId)
	if err := l.manager.updateLobby(l); err != nil {
		return err
	}
	return nil
}
