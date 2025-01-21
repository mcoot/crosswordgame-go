package types

import "github.com/hashicorp/go-uuid"

type PlayerId string

type Player struct {
	Id         PlayerId
	Name       string
	Registered bool
}

func NewEphemeralPlayer(name string) (*Player, error) {
	rawId, err := uuid.GenerateUUID()
	if err != nil {
		return nil, err
	}
	id := PlayerId(rawId)

	return &Player{
		Id:         id,
		Name:       name,
		Registered: false,
	}, nil
}
