package lobby

import "fmt"

type Manager struct {
	Lobbies map[LobbyID]*Lobby
}

func NewManager() Manager {
	return Manager{
		Lobbies: make(map[LobbyID]*Lobby),
	}
}

func (m *Manager) CreateLobby() *Lobby {
	lobby := NewLobby()
	m.Lobbies[lobby.ID] = lobby
	return lobby
}

func (m *Manager) GetLobby(id LobbyID) (*Lobby, error) {
	lobby, ok := m.Lobbies[id]
	if !ok {
		return nil, fmt.Errorf("lobby %s not found", id)
	}
	return lobby, nil
}

func (m *Manager) DeleteLobby(id LobbyID) {
	delete(m.Lobbies, id)
}
