package player

import "github.com/mcoot/crosswordgame-go/internal/store"

type Manager struct {
	store store.PlayerStore
}

func NewPlayerManager(store store.PlayerStore) *Manager {
	return &Manager{
		store: store,
	}
}
