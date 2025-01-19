package lobby

import (
	"github.com/hashicorp/go-uuid"
	"github.com/mcoot/crosswordgame-go/internal/lobby/types"
)

type Manager struct {
	Lobbies map[types.LobbyId]types.Lobby
}

func NewLobbyManager() *Manager {
	return &Manager{
		Lobbies: make(map[types.LobbyId]types.Lobby),
	}
}

func (m *Manager) CreateLobby(name string) (types.LobbyId, error) {
	rawId, err := uuid.GenerateUUID()
	if err != nil {
		return "", err
	}
	id := types.LobbyId(rawId)
	m.Lobbies[id] = types.Lobby{
		Name:    name,
		Players: make(map[types.PlayerId]types.Player),
	}
	return id, nil
}
