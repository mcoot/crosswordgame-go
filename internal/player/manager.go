package player

import (
	playertypes "github.com/mcoot/crosswordgame-go/internal/player/types"
	"github.com/mcoot/crosswordgame-go/internal/store"
)

type Manager struct {
	store store.PlayerStore
}

func NewPlayerManager(store store.PlayerStore) *Manager {
	return &Manager{
		store: store,
	}
}

func (m *Manager) LoginAsEphemeral(displayName string) (playertypes.PlayerId, error) {
	player, err := playertypes.NewEphemeralPlayer(displayName)
	if err != nil {
		return "", err
	}
	err = m.store.StorePlayer(player)
	if err != nil {
		return "", err
	}

	return player.Username, nil
}

func (m *Manager) LookupPlayer(playerId playertypes.PlayerId) (*playertypes.Player, error) {
	return m.store.RetrievePlayer(playerId)
}
