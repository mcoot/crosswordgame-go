package types

import (
	"fmt"
	"github.com/hashicorp/go-uuid"
	"github.com/mcoot/crosswordgame-go/internal/errors"
	"strings"
	"time"
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
	LastLogin   time.Time
}

func newPlayer(kind PlayerKind, username PlayerId, displayName string, lastLogin time.Time) *Player {
	return &Player{
		Kind:        kind,
		Username:    username,
		DisplayName: displayName,
		LastLogin:   lastLogin,
	}
}

func NewEphemeralPlayer(displayName string) (*Player, error) {
	rawId, err := uuid.GenerateUUID()
	if err != nil {
		return nil, err
	}
	id := PlayerId(fmt.Sprintf("%s%s", ephemeralPlayerPrefix, rawId))

	return newPlayer(PlayerKindEphemeral, id, displayName, time.Now()), nil
}

func NewRegisteredPlayer(username PlayerId, displayName string) (*Player, error) {
	if strings.HasPrefix(string(username), ephemeralPlayerPrefix) {
		return nil, &errors.InvalidInputError{
			ErrMessage: fmt.Sprintf("username cannot start with %s", ephemeralPlayerPrefix),
		}
	}

	// TODO: Creation shouldn't count as a login
	return newPlayer(PlayerKindRegistered, username, displayName, time.Now()), nil
}
