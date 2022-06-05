package lobby

import (
	"github.com/google/uuid"
)

type Manager struct {
	store *Store
}

func NewManager(store *Store) *Manager {
	return &Manager{
		store: store,
	}
}

func (m *Manager) CreateLobby() (*Lobby, error) {
	lobby := &Lobby{
		ID:      ID(uuid.New()),
		Members: map[int]struct{}{},
		manager: m,
	}
	err := m.store.Insert(lobby)
	if err != nil {
		return nil, err
	}
	return lobby, nil
}

func (m *Manager) GetLobby(id ID) (*Lobby, error) {
	lobby, err := m.store.Find(id)
	if err != nil {
		return nil, err
	}
	lobby.manager = m
	return lobby, nil
}

func (m *Manager) DeleteLobby(id ID) error {
	err := m.store.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (m *Manager) updateLobby(lobby *Lobby) error {
	err := m.store.Update(lobby)
	if err != nil {
		return err
	}
	return nil
}
