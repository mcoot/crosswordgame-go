package types

type PlayerId string

type Player struct {
	ID         PlayerId
	Name       string
	Registered bool
}
