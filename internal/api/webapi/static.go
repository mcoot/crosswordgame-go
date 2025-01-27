package webapi

import (
	"github.com/gorilla/mux"
	"net/http"
)

type StaticAssets struct{}

func NewStaticAssets() *StaticAssets {
	return &StaticAssets{}
}

func (s *StaticAssets) AttachToRouter(router *mux.Router) error {
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	return nil
}
