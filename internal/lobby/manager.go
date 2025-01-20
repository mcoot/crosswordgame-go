package lobby

import (
	"fmt"
	"github.com/hashicorp/go-uuid"
	"github.com/mcoot/crosswordgame-go/internal/errors"
	gametypes "github.com/mcoot/crosswordgame-go/internal/game/types"
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

func (m *Manager) RemovePlayerFromLobby(lobbyId types.LobbyId, playerId playertypes.PlayerId) error {
	lobby, err := m.store.RetrieveLobby(lobbyId)
	if err != nil {
		return err
	}
	for i, p := range lobby.Players {
		if p == playerId {
			lobby.Players = append(lobby.Players[:i], lobby.Players[i+1:]...)
			break
		}
	}
	return m.store.StoreLobby(lobbyId, lobby)
}

func (m *Manager) AttachGameToLobby(lobbyId types.LobbyId, gameId gametypes.GameId) error {
	lobby, err := m.store.RetrieveLobby(lobbyId)
	if err != nil {
		return err
	}

	if lobby.RunningGame != nil {
		return &errors.InvalidActionError{
			Action: "attach_game_to_lobby",
			Reason: fmt.Sprintf("lobby %s already has a running game, %s", lobbyId, lobby.RunningGame.GameId),
		}
	}

	lobby.RunningGame = &types.RunningGame{
		GameId: gameId,
	}
	return m.store.StoreLobby(lobbyId, lobby)
}

func (m *Manager) DetachGameFromLobby(lobbyId types.LobbyId) error {
	lobby, err := m.store.RetrieveLobby(lobbyId)
	if err != nil {
		return err
	}

	if lobby.RunningGame == nil {
		return &errors.InvalidActionError{
			Action: "detach_game_from_lobby",
			Reason: fmt.Sprintf("lobby %s does not have a running game", lobbyId),
		}
	}

	lobby.RunningGame = nil
	return m.store.StoreLobby(lobbyId, lobby)
}
