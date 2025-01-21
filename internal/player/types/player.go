package types

import (
	"github.com/hashicorp/go-uuid"
)

type PlayerKind string

const (
	PlayerKindRegistered = "registered"
	PlayerKindEphemeral  = "ephemeral"
)

type PlayerId string

type Player struct {
	Kind        PlayerKind
	Username    PlayerId
	DisplayName string
}

func NewPlayer(kind PlayerKind, username PlayerId, displayName string) *Player {
	return &Player{
		Kind:        kind,
		Username:    username,
		DisplayName: displayName,
	}
}

func NewEphemeralPlayer(displayName string) (*Player, error) {
	rawId, err := uuid.GenerateUUID()
	if err != nil {
		return nil, err
	}
	id := PlayerId(rawId)

	return NewPlayer(PlayerKindEphemeral, id, displayName), nil
}
