package cmd

import (
	"net/http"

	"github.com/gorilla/mux"
)

func SetupRoutes(router *mux.Router) {
    router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("frontend/static"))))
    router.HandleFunc("/", homeHandler).Methods("GET")
    router.HandleFunc("/chat", sendingMessageHandler).Methods("POST")
}