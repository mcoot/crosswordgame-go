package types

import (
	"fmt"
	"github.com/hashicorp/go-uuid"
	"github.com/mcoot/crosswordgame-go/internal/errors"
	"strings"
)

type PlayerKind string

const (
	PlayerKindRegistered = "registered"
	PlayerKindEphemeral  = "ephemeral"

	ephemeralPlayerPrefix = "ephemeral--"
)

type PlayerId string

type Player struct {
	Kind        PlayerKind
	Username    PlayerId
	DisplayName string
}

func newPlayer(kind PlayerKind, username PlayerId, displayName string) *Player {
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
	id := PlayerId(fmt.Sprintf("%s%s", ephemeralPlayerPrefix, rawId))

	return newPlayer(PlayerKindEphemeral, id, displayName), nil
}

func NewRegisteredPlayer(username PlayerId, displayName string) (*Player, error) {
	if strings.HasPrefix(string(username), ephemeralPlayerPrefix) {
		return nil, &errors.InvalidInputError{
			ErrMessage: fmt.Sprintf("username cannot start with %s", ephemeralPlayerPrefix),
		}
	}

	return newPlayer(PlayerKindRegistered, username, displayName), nil
}
