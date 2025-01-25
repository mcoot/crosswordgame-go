package webapi

import (
	"encoding/json"
	"fmt"
	"github.com/mcoot/crosswordgame-go/internal/api/utils"
	lobbytypes "github.com/mcoot/crosswordgame-go/internal/lobby/types"
	playertypes "github.com/mcoot/crosswordgame-go/internal/player/types"
	"net/http"
	"sync"
)

type sseEvent struct {
	Event string
	Data  interface{}
}

func (e sseEvent) Build() (string, error) {
	str, err := json.Marshal(e.Data)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("event: %s\r\ndata: %s\r\n\r\n", e.Event, str), nil
}

type refreshEvent struct {
	LobbyId          lobbytypes.LobbyId
	InitiatingPlayer playertypes.PlayerId
}

func (e refreshEvent) ToSSE() sseEvent {
	return sseEvent{
		Event: "refresh",
		Data:  e,
	}
}

type sseServer struct {
	running           bool
	refreshInput      chan refreshEvent
	activeConnections map[lobbytypes.LobbyId]map[playertypes.PlayerId][]chan refreshEvent
	mutex             sync.RWMutex
}

func newSSEServer() *sseServer {
	return &sseServer{
		refreshInput:      make(chan refreshEvent),
		activeConnections: make(map[lobbytypes.LobbyId]map[playertypes.PlayerId][]chan refreshEvent),
		mutex:             sync.RWMutex{},
	}
}

func (s *sseServer) HandleRequest(w http.ResponseWriter, r *http.Request, player *playertypes.Player) {
	lobbyId := utils.GetLobbyIdPathParam(r)

	// Set http headers required for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// You may need this locally for CORS requests
	//w.Header().Set("Access-Control-Allow-Origin", "*")

	eventChan := s.newConnection(lobbyId, player.Username)

	rc := http.NewResponseController(w)
	for {
		select {
		case <-r.Context().Done():
			s.dropConnection(lobbyId, player.Username, eventChan)
			return
		case <-eventChan:
			evt, err := refreshEvent{}.ToSSE().Build()
			if err != nil {
				return
			}
			_, err = w.Write([]byte(evt))
			if err != nil {
				return
			}
			err = rc.Flush()
			if err != nil {
				return
			}
		}
	}

}

func (s *sseServer) Start() {
	if s.running {
		return
	}
	s.running = true
	defer func() {
		s.running = false
	}()
	// TODO: context for graceful shutdown
	go func() {
		for e := range s.refreshInput {
			s.broadcastRefreshExceptToInitiator(e.LobbyId, e.InitiatingPlayer)
		}
	}()
}

func (s *sseServer) SendRefresh(lobbyId lobbytypes.LobbyId, initiatingPlayer playertypes.PlayerId) {
	s.refreshInput <- refreshEvent{
		LobbyId:          lobbyId,
		InitiatingPlayer: initiatingPlayer,
	}
}

func (s *sseServer) newConnection(lobbyId lobbytypes.LobbyId, playerId playertypes.PlayerId) chan refreshEvent {
	channel := make(chan refreshEvent)
	s.mutex.Lock()
	defer s.mutex.Unlock()
	channelsForLobby, ok := s.activeConnections[lobbyId]
	if !ok {
		channelsForLobby = make(map[playertypes.PlayerId][]chan refreshEvent)
	}

	channelsForPlayer, ok := channelsForLobby[playerId]
	if !ok {
		channelsForPlayer = make([]chan refreshEvent, 0)
	}
	channelsForPlayer = append(channelsForPlayer, channel)
	channelsForLobby[playerId] = channelsForPlayer
	s.activeConnections[lobbyId] = channelsForLobby
	return channel
}

func (s *sseServer) dropConnection(
	lobbyId lobbytypes.LobbyId,
	playerId playertypes.PlayerId,
	channel chan refreshEvent,
) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	channelsForLobby, ok := s.activeConnections[lobbyId]
	if !ok {
		return
	}

	channelsForPlayer, ok := channelsForLobby[playerId]
	if !ok {
		return
	}
	for i, c := range channelsForPlayer {
		if c == channel {
			channelsForPlayer = append(channelsForPlayer[:i], channelsForPlayer[i+1:]...)
			close(c)
			break
		}
	}
	channelsForLobby[playerId] = channelsForPlayer
	s.activeConnections[lobbyId] = channelsForLobby
}

func (s *sseServer) broadcastRefreshExceptToInitiator(lobbyId lobbytypes.LobbyId, initiator playertypes.PlayerId) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	channelsForLobby, ok := s.activeConnections[lobbyId]
	if !ok {
		return
	}

	for playerId, channels := range channelsForLobby {
		if playerId != initiator {
			for _, c := range channels {
				c <- refreshEvent{
					LobbyId:          lobbyId,
					InitiatingPlayer: initiator,
				}
			}
		}
	}
}
