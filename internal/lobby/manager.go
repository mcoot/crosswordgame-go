package lobby

import (
	"github.com/hashicorp/go-uuid"
	"github.com/mcoot/crosswordgame-go/internal/lobby/types"
	playertypes "github.com/mcoot/crosswordgame-go/internal/player/types"
	"github.com/mcoot/crosswordgame-go/internal/store"
)

type Manager struct {
	store store.LobbyStore
}

func NewLobbyManager(store store.LobbyStore) *Manager {
	return &Manager{
		store: store,
	}
}

func (m *Manager) CreateLobby(name string) (types.LobbyId, error) {
	rawId, err := uuid.GenerateUUID()
	if err != nil {
		return "", err
	}
	id := types.LobbyId(rawId)
	err = m.store.StoreLobby(id, &types.Lobby{
		Name:        name,
		Players:     make([]playertypes.PlayerId, 0),
		RunningGame: nil,
	})
	if err != nil {
		return "", err
	}
	return id, nil
}

func (m *Manager) GetLobbyState(id types.LobbyId) (*types.Lobby, error) {
	lobby, err := m.store.RetrieveLobby(id)
	if err != nil {
		return nil, err
	}
	return lobby, nil
}

func (m *Manager) AddPlayerToLobby(lobbyId types.LobbyId, playerId playertypes.PlayerId) error {
	lobby, err := m.store.RetrieveLobby(lobbyId)
	if err != nil {
		return err
	}
	lobby.Players = append(lobby.Players, playerId)
	return m.store.StoreLobby(lobbyId, lobby)
}
