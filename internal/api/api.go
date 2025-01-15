package api

import "net/http"

type CrosswordGameAPI struct{}

func (c *CrosswordGameAPI) AttachToMux(h *http.ServeMux) {
	h.Handle("GET /health", http.HandlerFunc(c.Healthcheck))
	h.Handle("POST /api/v1/game", http.HandlerFunc(c.CreateGame))
	h.Handle("GET /api/v1/game/{gameId}", http.HandlerFunc(c.GetGameState))
	h.Handle("GET /api/v1/game/{gameId}/player/{playerId}", http.HandlerFunc(c.GetPlayerState))
	h.Handle("POST /api/v1/game/{gameId}/player/{playerId}/announce", http.HandlerFunc(c.SubmitAnnouncement))
	h.Handle("POST /api/v1/game/{gameId}/player/{playerId}/place", http.HandlerFunc(c.SubmitPlacement))
	h.Handle("GET /api/v1/game/{gameId}/player/{playerId}/score", http.HandlerFunc(c.GetPlayerScore))
}

func (c *CrosswordGameAPI) Healthcheck(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("OK"))
}

func (c *CrosswordGameAPI) CreateGame(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(500)
}

func (c *CrosswordGameAPI) GetGameState(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(500)
}

func (c *CrosswordGameAPI) GetPlayerState(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(500)
}

func (c *CrosswordGameAPI) SubmitAnnouncement(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(500)
}

func (c *CrosswordGameAPI) SubmitPlacement(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(500)
}

func (c *CrosswordGameAPI) GetPlayerScore(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(500)
}
